package worker

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusOffline Status = iota
	StatusActive
	StatusBusy
)

// Worker Aggregate
type Worker struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Nome          string
	Status        Status
	LastHeartbeat *time.Time
	Capacidades   string // Json
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type Repository interface {
	Save(w *Worker) error
	FindByName(nome string) (*Worker, error)
	UpdateHeartbeat(id uuid.UUID) error
}
