package job

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJobEntity_StatusConstants(t *testing.T) {
	if StatusQueued != 0 {
		t.Errorf("expected StatusQueued to be 0, got %v", StatusQueued)
	}
	if StatusProcessing != 1 {
		t.Errorf("expected StatusProcessing to be 1, got %v", StatusProcessing)
	}
	if StatusCompleted != 2 {
		t.Errorf("expected StatusCompleted to be 2, got %v", StatusCompleted)
	}
	if StatusFailed != 3 {
		t.Errorf("expected StatusFailed to be 3, got %v", StatusFailed)
	}
}

func TestJobEntity_Struct(t *testing.T) {
	now := time.Now()
	params := "{}"
	j := JobProcessamento{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		VideoID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		WorkerID:   nil,
		Tipo:       1,
		Status:     StatusQueued,
		StartedAt:  &now,
		FinishedAt: nil,
		Params:     &params,
	}

	if j.Tipo != 1 {
		t.Errorf("expected Tipo %v, got %v", 1, j.Tipo)
	}
	if j.Status != StatusQueued {
		t.Errorf("expected Status %v, got %v", StatusQueued, j.Status)
	}
	if j.StartedAt != &now {
		t.Errorf("expected StartedAt to be set")
	}
	if j.FinishedAt != nil {
		t.Errorf("expected FinishedAt nil, got %v", j.FinishedAt)
	}
}
