package frame

import (
	"time"

	"github.com/google/uuid"
)

// Frame Aggregate
type Frame struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	VideoID     uuid.UUID
	FrameNo     int64
	TimestampMs int64
	Features    string // Representado como string JSON para simplificar
	URI         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Repository define a interface de persistência para Frame
type Repository interface {
	Save(frame *Frame) error
	FindByVideoID(videoID uuid.UUID) ([]*Frame, error)
}
