package video

import (
	"time"
	"github.com/google/uuid"
)

// Status define os possíveis estados de um vídeo.
type Status int16

const (
	StatusPending    Status = 0
	StatusProcessing Status = 1
	StatusCompleted  Status = 2
	StatusFailed     Status = 3
)

// Video representa a entidade 'video' do banco de dados.
type Video struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Title     string     `json:"title"`
	VideoPath string     `json:"video_path"` // Caminho do objeto no storage
	Status    Status     `json:"status"`
	Meta      *string    `json:"meta,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
