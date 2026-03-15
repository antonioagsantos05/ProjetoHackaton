package job

import (
	"time"
	"github.com/google/uuid"
)

// Status define os possíveis estados de um job.
type Status int16

const (
	StatusQueued     Status = 0
	StatusProcessing Status = 1
	StatusCompleted  Status = 2
	StatusFailed     Status = 3
)

// Job is an alias for JobProcessamento for API compatibility.
type Job = JobProcessamento

// JobProcessamento representa a entidade 'job_processamento' do banco de dados.
type JobProcessamento struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	VideoID     uuid.UUID  `json:"video_id"`
	WorkerID    *uuid.UUID `json:"worker_id,omitempty"`
	Tipo        int16      `json:"tipo"`
	Status      Status     `json:"status"`
	VideoPath   string     `json:"video_path"` // Campo restaurado
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	Params      *string    `json:"params,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
