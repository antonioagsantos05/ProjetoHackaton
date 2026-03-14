package session

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSessionEntity_Struct(t *testing.T) {
	now := time.Now()
	s := Session{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Token:     "ey.jwt.token",
		ExpiresAt: now.Add(1 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.Token != "ey.jwt.token" {
		t.Errorf("expected Token ey.jwt.token, got %v", s.Token)
	}
	if s.ExpiresAt.Before(now) {
		t.Errorf("expected ExpiresAt to be in the future")
	}
}