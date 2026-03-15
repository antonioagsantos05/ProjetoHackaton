package video

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/minio/minio-go/v7"
)

// DownloadUseCase define a interface para o caso de uso de download.
type DownloadUseCase interface {
	Execute(jobID string) (string, error)
}

// downloadUseCase é a implementação do DownloadUseCase.
type downloadUseCase struct {
	db          *sql.DB
	minioClient *minio.Client
}

// NewDownloadUseCase cria uma nova instância do downloadUseCase.
func NewDownloadUseCase(db *sql.DB, minioClient *minio.Client) DownloadUseCase {
	return &downloadUseCase{db: db, minioClient: minioClient}
}

// Execute orquestra o processo de obtenção do link de download do ZIP.
func (uc *downloadUseCase) Execute(jobID string) (string, error) {
	// 1. Buscar o status do job e a URI do arquivo ZIP
	var jobStatus job.Status
	// Usando sql.NullString pois o LEFT JOIN pode retornar nulo se o job ainda não terminou
	var zipURINull sql.NullString

	query := `
		SELECT jp.status, az.uri
		FROM job_processamento jp
		LEFT JOIN arquivo_zip az ON jp.id = az.job_id
		WHERE jp.id = $1
	`
	err := uc.db.QueryRow(query, jobID).Scan(&jobStatus, &zipURINull)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("job com ID '%s' não encontrado", jobID)
		}
		return "", fmt.Errorf("erro ao buscar job e arquivo zip no banco de dados: %w", err)
	}

	// 2. Verificar o status do job
	if jobStatus != job.StatusCompleted {
		return "", fmt.Errorf("o processamento do job '%s' ainda não foi concluído. Status atual: %d", jobID, jobStatus)
	}

	if !zipURINull.Valid || zipURINull.String == "" {
		return "", fmt.Errorf("o processamento foi concluído, mas o arquivo ZIP não foi encontrado no banco de dados")
	}

	zipURI := zipURINull.String

	// 3. Gerar URL pré-assinada do Minio para o arquivo ZIP.
	expiration := time.Hour * 1

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"resultado.zip\"")

	presignedURL, err := uc.minioClient.PresignedGetObject(context.Background(), MinioBucket, zipURI, expiration, reqParams)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar URL pré-assinada do Minio: %w", err)
	}

	return presignedURL.String(), nil
}
