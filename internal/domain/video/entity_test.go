package video

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestVideoEntity_StatusConstants(t *testing.T) {
	if StatusPending != 0 {
		t.Errorf("expected StatusPending to be 0, got %v", StatusPending)
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

func TestVideoEntity_Struct(t *testing.T) {
	now := time.Now()
	v := Video{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Title:     "My Video",
		Status:    StatusPending,
		Meta:      "{}",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if v.Title != "My Video" {
		t.Errorf("expected Title %v, got %v", "My Video", v.Title)
	}
	if v.Status != StatusPending {
		t.Errorf("expected Status %v, got %v", StatusPending, v.Status)
	}
	if v.CreatedAt != now {
		t.Errorf("expected CreatedAt %v, got %v", now, v.CreatedAt)
	}
}
