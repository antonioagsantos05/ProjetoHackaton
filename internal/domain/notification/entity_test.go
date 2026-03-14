package notification

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNotificationEntity_Struct(t *testing.T) {
	now := time.Now()
	n := Notificacao{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Tipo:      TipoSucesso,
		Payload:   `{"msg": "Processamento Concluido"}`,
		Lida:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if n.Tipo != TipoSucesso {
		t.Errorf("expected TipoSucesso, got %v", n.Tipo)
	}
	if n.Lida {
		t.Errorf("expected Lida to be false")
	}
}