package postgres

import (
	"database/sql"
	"errors"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/google/uuid"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *domainJob.Job) error {
	// Exemplo de como persistir no PostgreSQL com tenant_id
	// A tabela "job_processamento" requeriria a query de insert.
	return nil // mock
}

func (r *JobRepository) Update(job *domainJob.Job) error {
	return errors.New("not implemented")
}

func (r *JobRepository) FindByID(id uuid.UUID) (*domainJob.Job, error) {
	return nil, errors.New("not implemented")
}
