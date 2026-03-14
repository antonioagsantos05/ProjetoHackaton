package policy

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPolicyEntity_Struct(t *testing.T) {
	now := time.Now()
	p := Retencao{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Scope:     ScopeTenant,
		Dias:      30,
		Ativo:     true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if p.Scope != ScopeTenant {
		t.Errorf("expected ScopeTenant, got %v", p.Scope)
	}
	if p.Dias != 30 {
		t.Errorf("expected Dias 30, got %v", p.Dias)
	}
	if !p.Ativo {
		t.Errorf("expected Ativo to be true")
	}
}