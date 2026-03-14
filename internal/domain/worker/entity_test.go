package worker

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestWorkerEntity_Struct(t *testing.T) {
	now := time.Now()
	w := Worker{
		ID:            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:      uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Nome:          "worker-gpu-1",
		Status:        StatusActive,
		LastHeartbeat: &now,
		Capacidades:   `{"gpu": true, "max_jobs": 5}`,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if w.Status != StatusActive {
		t.Errorf("expected StatusActive, got %v", w.Status)
	}
	if w.Nome != "worker-gpu-1" {
		t.Errorf("expected worker-gpu-1, got %v", w.Nome)
	}
}