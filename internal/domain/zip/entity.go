package zip

import (
	"time"

	"github.com/google/uuid"
)

// ArquivoZip Aggregate
type ArquivoZip struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	JobID     uuid.UUID
	URI       string
	Checksum  string
	Tamanho   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Repository define a interface de persistência para Zip
type Repository interface {
	Save(zip *ArquivoZip) error
	FindByJobID(jobID uuid.UUID) (*ArquivoZip, error)
}
