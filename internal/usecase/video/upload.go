package video

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/streadway/amqp"
)

const (
	VideoProcessingQueue = "video_processing"
	MinioBucket          = "videos"
)

// UploadUseCase define a interface para o caso de uso de upload.
type UploadUseCase interface {
	Execute(localFilePath string, userID, tenantID uuid.UUID) (*job.JobProcessamento, error)
}

// uploadUseCase é a implementação do UploadUseCase.
type uploadUseCase struct {
	db          *sql.DB
	ch          *amqp.Channel
	minioClient *minio.Client
}

// NewUploadUseCase cria uma nova instância do uploadUseCase.
func NewUploadUseCase(db *sql.DB, ch *amqp.Channel, minioClient *minio.Client) UploadUseCase {
	return &uploadUseCase{db: db, ch: ch, minioClient: minioClient}
}

// Execute orquestra o processo de upload.
func (uc *uploadUseCase) Execute(localFilePath string, userID, tenantID uuid.UUID) (*job.JobProcessamento, error) {
	ctx := context.Background()
	objectName := uuid.New().String()

	// 1. Fazer upload do arquivo para o Minio
	_, err := uc.minioClient.FPutObject(ctx, MinioBucket, objectName, localFilePath, minio.PutObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do arquivo para o Minio: %w", err)
	}
	log.Printf("Arquivo '%s' enviado para o Minio como '%s'", localFilePath, objectName)
	defer os.Remove(localFilePath)

	// 2. Criar e salvar a entidade Video
	newVideo := &video.Video{
		ID:        uuid.New(),
		TenantID:  tenantID,
		UserID:    userID,
		Title:     filepath.Base(localFilePath),
		VideoPath: objectName,
		Status:    video.StatusPending,
		CreatedAt: time.Now(),
	}
	videoQuery := `INSERT INTO video (id, tenant_id, usuario_id, titulo, video_path, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = uc.db.Exec(videoQuery, newVideo.ID, newVideo.TenantID, newVideo.UserID, newVideo.Title, newVideo.VideoPath, newVideo.Status, newVideo.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir vídeo no banco de dados: %w", err)
	}
	log.Printf("Vídeo %s salvo no banco de dados", newVideo.ID)

	// 3. Criar e salvar a entidade JobProcessamento
	newJob := &job.JobProcessamento{
		ID:        uuid.New(),
		TenantID:  tenantID,
		VideoID:   newVideo.ID,
		Tipo:      1,
		Status:    job.StatusQueued,
		CreatedAt: time.Now(),
	}
	jobQuery := `INSERT INTO job_processamento (id, tenant_id, video_id, tipo, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = uc.db.Exec(jobQuery, newJob.ID, newJob.TenantID, newJob.VideoID, newJob.Tipo, newJob.Status, newJob.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir job no banco de dados: %w", err)
	}
	log.Printf("Job %s salvo no banco de dados com status QUEUED", newJob.ID)

	// 4. Declarar a fila (garante que ela exista antes de publicar)
	_, err = uc.ch.QueueDeclare(
		VideoProcessingQueue, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao declarar a fila no RabbitMQ: %w", err)
	}

	// 5. Publicar mensagem no RabbitMQ
	jobJSON, err := json.Marshal(newJob)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar job para JSON: %w", err)
	}

	err = uc.ch.Publish(
		"",                   // exchange
		VideoProcessingQueue, // routing key (nome da fila)
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jobJSON,
		})
	if err != nil {
		return nil, fmt.Errorf("erro ao publicar mensagem no RabbitMQ: %w", err)
	}
	log.Printf("Job %s publicado na fila '%s'", newJob.ID, VideoProcessingQueue)

	return newJob, nil
}
