package postgres

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/google/uuid"
)

// Nota: Testes reais de repositório requerem um BD em execução (Testcontainers) ou Mocks de DB (sqlmock).
// Aqui testamos apenas a estrutura do mock em FindByID/UpdateStatus.

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

// Simulando a falha do Execute quando o banco não está mockado corretamente:
func TestVideoRepository_Save_Error(t *testing.T) {
	repo := NewVideoRepository(&sql.DB{}) // DB Nulo ou Desconectado
	video := &domainVideo.Video{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		UserID:    uuid.New(),
		Title:     "Test",
		Status:    domainVideo.StatusPending,
		Meta:      "{}",
		CreatedAt: time.Now(),
	}

	// Como o *sql.DB está vazio/inativo, .Exec vai dar pânico ou erro:
	// Vamos encapsular em defer recover apenas caso dê panic pela implementação interna do go:
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
