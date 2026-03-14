package session

import (
	"time"

	"github.com/google/uuid"
)

// Session Aggregate Root
type Session struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Repository define a interface de persistência para Sessao
type Repository interface {
	Save(session *Session) error
	FindByToken(token string) (*Session, error)
	Delete(id uuid.UUID) error
}
