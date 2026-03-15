package job

import (
	"database/sql"
	"fmt"

	"github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/fiap-x/video-processor/internal/domain/video" // Importar o domain de video
)

// StatusResponse é a estrutura de dados que será retornada, combinando informações.
type StatusResponse struct {
	JobInfo   job.JobProcessamento `json:"job_info"`
	VideoInfo video.Video            `json:"video_info"`
}

// StatusUseCase define a interface para o caso de uso de verificação de status de job.
type StatusUseCase interface {
	Execute(jobID string) (*StatusResponse, error)
}

// statusUseCase é a implementação de StatusUseCase.
type statusUseCase struct {
	db *sql.DB
}

// NewStatusUseCase cria uma nova instância de statusUseCase.
func NewStatusUseCase(db *sql.DB) StatusUseCase {
	return &statusUseCase{db: db}
}

// Execute recupera um job e as informações do vídeo associado.
func (uc *statusUseCase) Execute(jobID string) (*StatusResponse, error) {
	var j job.JobProcessamento
	var v video.Video

	query := `
		SELECT
			jp.id, jp.tenant_id, jp.video_id, jp.tipo, jp.status, jp.created_at,
			v.id, v.tenant_id, v.usuario_id, v.titulo, v.video_path, v.status, v.created_at
		FROM job_processamento jp
		JOIN video v ON jp.video_id = v.id
		WHERE jp.id = $1
	`
	err := uc.db.QueryRow(query, jobID).Scan(
		&j.ID, &j.TenantID, &j.VideoID, &j.Tipo, &j.Status, &j.CreatedAt,
		&v.ID, &v.TenantID, &v.UserID, &v.Title, &v.VideoPath, &v.Status, &v.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job com ID '%s' não encontrado", jobID)
		}
		return nil, fmt.Errorf("erro ao buscar job e vídeo no banco de dados: %w", err)
	}

	response := &StatusResponse{
		JobInfo:   j,
		VideoInfo: v,
	}

	return response, nil
}
