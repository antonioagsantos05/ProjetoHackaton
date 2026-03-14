package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserEntity_Struct(t *testing.T) {
	now := time.Now()
	u := User{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Email:     "test@fiap.com.br",
		Name:      "Test User",
		HashPass:  "hash123",
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if u.Email != "test@fiap.com.br" {
		t.Errorf("expected email %v, got %v", "test@fiap.com.br", u.Email)
	}
	if u.Status != 1 {
		t.Errorf("expected status %v, got %v", 1, u.Status)
	}
	if u.CreatedAt != now {
		t.Errorf("expected CreatedAt %v, got %v", now, u.CreatedAt)
	}
}
