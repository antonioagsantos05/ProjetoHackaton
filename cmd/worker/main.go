package main

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/fiap-x/video-processor/internal/infrastructure/database/postgres"
	"github.com/fiap-x/video-processor/internal/infrastructure/messaging/rabbitmq"
	"github.com/fiap-x/video-processor/internal/infrastructure/storage/minio"
	"github.com/google/uuid"
	minioGo "github.com/minio/minio-go/v7"
	"github.com/streadway/amqp"
)

const (
	VideoProcessingQueue = "video_processing"
	MinioBucket          = "videos"
)

func main() {
	// --- Conexões ---
	db, err := postgres.NewDBConnection()
	if err != nil {
		log.Fatalf("Falha ao conectar no DB: %v", err)
	}
	defer db.Close()

	rabbitConn, rabbitCh, err := rabbitmq.NewRabbitMQConnection()
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	minioClient, err := minio.NewMinioClient()
	if err != nil {
		log.Fatalf("Falha ao conectar no Minio: %v", err)
	}

	// --- Declara a fila (garante que ela exista) ---
	_, err = rabbitCh.QueueDeclare(
		VideoProcessingQueue, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Falha ao declarar a fila: %v", err)
	}

	// --- Configuração do Consumidor RabbitMQ ---
	msgs, err := rabbitCh.Consume(
		VideoProcessingQueue, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Fatalf("Falha ao registrar consumidor: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Recebida nova mensagem: %s", d.Body)
			processMessage(d, db, minioClient)
		}
	}()

	log.Printf(" [*] Worker aguardando por mensagens. Para sair, pressione CTRL+C")
	<-forever
}

func processMessage(d amqp.Delivery, db *sql.DB, minioClient *minioGo.Client) {
	var j job.JobProcessamento
	if err := json.Unmarshal(d.Body, &j); err != nil {
		log.Printf("Erro ao decodificar JSON do job: %s", err)
		return
	}

	log.Printf("Processando job %s para o vídeo %s...", j.ID, j.VideoID)

	// 1. Atualiza o status para PROCESSING
	updateStatus(db, j.ID, j.VideoID, job.StatusProcessing)

	// 2. Processa o vídeo real (baixa, zipa, sobe)
	err := processVideoReal(j, db, minioClient)
	if err != nil {
		log.Printf("Erro ao processar vídeo do job %s: %v", j.ID, err)
		updateStatus(db, j.ID, j.VideoID, job.StatusFailed)
		return
	}

	// 3. Atualiza o status para COMPLETED
	updateStatus(db, j.ID, j.VideoID, job.StatusCompleted)

	log.Printf("Job %s finalizado com sucesso!", j.ID)
}

func processVideoReal(j job.JobProcessamento, db *sql.DB, minioClient *minioGo.Client) error {
	ctx := context.Background()

	// 1. Buscar o video_path no BD (como o arquivo está salvo no minio)
	var videoPath string
	err := db.QueryRow(`SELECT video_path FROM video WHERE id = $1`, j.VideoID).Scan(&videoPath)
	if err != nil {
		return fmt.Errorf("erro ao buscar video_path: %w", err)
	}

	// 2. Fazer o download do vídeo do Minio para o worker
	localVideoPath := filepath.Join("/tmp", videoPath)
	err = minioClient.FGetObject(ctx, MinioBucket, videoPath, localVideoPath, minioGo.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("erro ao baixar video do minio: %w", err)
	}
	defer os.Remove(localVideoPath)

	// 3. Criar o arquivo ZIP contendo o vídeo
	zipObjectName := uuid.New().String() + ".zip"
	localZipPath := filepath.Join("/tmp", zipObjectName)

	err = createZip(localZipPath, localVideoPath, "video_processado.mp4")
	if err != nil {
		return fmt.Errorf("erro ao criar zip: %w", err)
	}
	defer os.Remove(localZipPath)

	// 4. Fazer upload do arquivo ZIP pro Minio
	_, err = minioClient.FPutObject(ctx, MinioBucket, zipObjectName, localZipPath, minioGo.PutObjectOptions{
		ContentType: "application/zip",
	})
	if err != nil {
		return fmt.Errorf("erro ao fazer upload do zip pro minio: %w", err)
	}

	// 5. Salvar registro na tabela arquivo_zip
	zipID := uuid.New()
	query := `INSERT INTO arquivo_zip (id, tenant_id, job_id, uri, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, zipID, j.TenantID, j.ID, zipObjectName, time.Now())
	if err != nil {
		return fmt.Errorf("erro ao salvar registro arquivo_zip: %w", err)
	}

	return nil
}

// createZip cria um arquivo .zip contendo um único arquivo com o nome especificado
func createZip(zipPath, sourceFile, nameInZip string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	writer, err := archive.Create(nameInZip)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, source)
	return err
}

func updateStatus(db *sql.DB, jobID, videoID uuid.UUID, newStatus job.Status) {
	// Atualiza o status do job
	jobQuery := `UPDATE job_processamento SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := db.Exec(jobQuery, newStatus, time.Now(), jobID)
	if err != nil {
		log.Printf("Erro ao atualizar status do job %s para %d: %v", jobID, newStatus, err)
	} else {
		log.Printf("Status do job %s atualizado para %d", jobID, newStatus)
	}

	// Se o job foi concluído ou falhou, atualiza também o status do vídeo
	if newStatus == job.StatusCompleted || newStatus == job.StatusFailed {
		videoStatus := video.StatusCompleted
		if newStatus == job.StatusFailed {
			videoStatus = video.StatusFailed
		}

		videoQuery := `UPDATE video SET status = $1, updated_at = $2 WHERE id = $3`
		_, err := db.Exec(videoQuery, videoStatus, time.Now(), videoID)
		if err != nil {
			log.Printf("Erro ao atualizar status do vídeo %s para %d: %v", videoID, videoStatus, err)
		} else {
			log.Printf("Status do vídeo %s atualizado para %d", videoID, videoStatus)
		}
	}
}
