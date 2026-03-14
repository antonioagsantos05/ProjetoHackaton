package frame

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestFrameEntity_Struct(t *testing.T) {
	now := time.Now()
	f := Frame{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		VideoID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		FrameNo:     10,
		TimestampMs: 1500,
		Features:    "{}",
		URI:         "s3://bucket/frame10.jpg",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if f.FrameNo != 10 {
		t.Errorf("expected FrameNo 10, got %v", f.FrameNo)
	}
	if f.URI != "s3://bucket/frame10.jpg" {
		t.Errorf("expected URI %v, got %v", "s3://bucket/frame10.jpg", f.URI)
	}
}