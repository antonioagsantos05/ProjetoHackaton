package job

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusQueued Status = iota
	StatusRunning
	StatusDone
	StatusError
)

// Job Aggregate Root para o processamento assíncrono
type Job struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	VideoID    uuid.UUID
	WorkerID   *uuid.UUID
	Type       int
	Status     Status
	Params     string // Representado como JSON
	StartedAt  *time.Time
	FinishedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time // Suporte a soft delete conforme DDL
}

// Repository define a interface de persistência para Job
type Repository interface {
	Create(job *Job) error
	Update(job *Job) error
	FindByID(id uuid.UUID) (*Job, error)
}

// MessageQueue define a interface para envio de jobs para as filas
type MessageQueue interface {
	PublishJob(job *Job) error
}
