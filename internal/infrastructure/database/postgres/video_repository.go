package postgres

import (
	"database/sql"
	"errors"

	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/google/uuid"
)

type VideoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Save(video *domainVideo.Video) error {
	query := `
		INSERT INTO video (id, tenant_id, usuario_id, titulo, status, meta, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		video.ID, video.TenantID, video.UserID, video.Title,
		video.Status, video.Meta, video.CreatedAt, video.UpdatedAt)
	return err
}

func (r *VideoRepository) FindByID(id uuid.UUID) (*domainVideo.Video, error) {
	// Implementação simulada para demonstração
	return nil, errors.New("not implemented")
}

func (r *VideoRepository) UpdateStatus(id uuid.UUID, status domainVideo.Status) error {
	// Implementação simulada para demonstração
	return errors.New("not implemented")
}
