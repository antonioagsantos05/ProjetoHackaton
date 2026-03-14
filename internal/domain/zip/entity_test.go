package zip

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestZipEntity_Struct(t *testing.T) {
	now := time.Now()
	z := ArquivoZip{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		JobID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		URI:       "s3://bucket/job_123.zip",
		Checksum:  "hash123",
		Tamanho:   1024,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if z.Tamanho != 1024 {
		t.Errorf("expected Tamanho 1024, got %v", z.Tamanho)
	}
	if z.Checksum != "hash123" {
		t.Errorf("expected Checksum hash123, got %v", z.Checksum)
	}
}