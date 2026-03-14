package postgres

import (
	"database/sql"
	"errors"
	"testing"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/google/uuid"
)

func TestJobRepository_Create_NotImplemented(t *testing.T) {
	repo := NewJobRepository(&sql.DB{})

	err := repo.Create(&domainJob.Job{})
	if err != nil {
		t.Errorf("expected mock nil error, got %v", err)
	}
}

func TestJobRepository_Update_NotImplemented(t *testing.T) {
	repo := NewJobRepository(&sql.DB{})

	err := repo.Update(&domainJob.Job{})
	if err == nil {
		t.Errorf("expected not implemented error, got nil")
	}
}

func TestJobRepository_FindByID_NotImplemented(t *testing.T) {
	repo := NewJobRepository(&sql.DB{})

	_, err := repo.FindByID(uuid.New())
	if err == nil {
		t.Errorf("expected not implemented error, got nil")
	}
}
