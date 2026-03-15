package postgres

import (
	"database/sql"
	"testing"
	"time"

	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/google/uuid"
)

func TestVideoRepository_FindByID_NotImplemented(t *testing.T) {
	repo := NewVideoRepository(&sql.DB{})

	_, err := repo.FindByID(uuid.New())
	if err == nil {
		t.Errorf("expected not implemented error, got nil")
	}
}

func TestVideoRepository_UpdateStatus_NotImplemented(t *testing.T) {
	repo := NewVideoRepository(&sql.DB{})

	err := repo.UpdateStatus(uuid.New(), domainVideo.StatusPending)
	if err == nil {
		t.Errorf("expected not implemented error, got nil")
	}
}

func TestVideoRepository_Save_Error(t *testing.T) {
	repo := NewVideoRepository(&sql.DB{})
	meta := "{}"
	video := &domainVideo.Video{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		UserID:    uuid.New(),
		Title:     "Test",
		Status:    domainVideo.StatusPending,
		Meta:      &meta,
		CreatedAt: time.Now(),
	}

	defer func() {
		if r := recover(); r != nil {
			// Panic esperado por causa do sql.DB{} não inicializado
		}
	}()

	err := repo.Save(video)
	if err == nil {
		t.Errorf("expected error or panic, got nil")
	}
}
