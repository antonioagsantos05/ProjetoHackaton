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
	if StatusRunning != 1 {
		t.Errorf("expected StatusRunning to be 1, got %v", StatusRunning)
	}
	if StatusDone != 2 {
		t.Errorf("expected StatusDone to be 2, got %v", StatusDone)
	}
	if StatusError != 3 {
		t.Errorf("expected StatusError to be 3, got %v", StatusError)
	}
}

func TestJobEntity_Struct(t *testing.T) {
	now := time.Now()
	j := Job{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		VideoID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		WorkerID:   nil,
		Type:       1,
		Status:     StatusQueued,
		StartedAt:  &now,
		FinishedAt: nil,
		Params:     "{}",
	}

	if j.Type != 1 {
		t.Errorf("expected Type %v, got %v", 1, j.Type)
	}
	if j.Status != StatusQueued {
		t.Errorf("expected Status %v, got %v", StatusQueued, j.Status)
	}
	if j.StartedAt != &now {
		t.Errorf("expected StartedAt %v, got %v", &now, j.StartedAt)
	}
	if j.FinishedAt != nil {
		t.Errorf("expected FinishedAt nil, got %v", j.FinishedAt)
	}
}
