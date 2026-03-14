package notification

import (
	"time"

	"github.com/google/uuid"
)

type Tipo int

const (
	TipoSucesso Tipo = iota
	TipoErro
	TipoAviso
)

// Notificacao Aggregate
type Notificacao struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Tipo      Tipo
	Payload   string // Representado como string JSON para simplificar
	Lida      bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Repository define a interface de persistência para Notificacao
type Repository interface {
	Save(notificacao *Notificacao) error
	FindByUserID(userID uuid.UUID) ([]*Notificacao, error)
	MarkAsRead(id uuid.UUID) error
}
