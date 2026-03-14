package video

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusPending Status = iota
	StatusProcessing
	StatusCompleted
	StatusFailed
)

// Video Aggregate Root
type Video struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Title     string
	Status    Status
	Meta      string // Representado como string JSON para simplificar
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // Suporte a soft delete conforme DDL
}

// Repository define a interface de persistência para Video
type Repository interface {
	Save(video *Video) error
	FindByID(id uuid.UUID) (*Video, error)
	UpdateStatus(id uuid.UUID, status Status) error
}
