package user

import (
	"time"
	"github.com/google/uuid"
)

// User representa a entidade 'usuario' do banco de dados.
type User struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Email     string     `json:"email"`
	Nome      string     `json:"nome"`
	HashSenha string     `json:"-"` // O hash da senha nunca deve ser exposto
	Status    int16      `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
