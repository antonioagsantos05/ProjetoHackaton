package policy

import (
	"time"

	"github.com/google/uuid"
)

type Scope int

const (
	ScopeGlobal Scope = iota
	ScopeTenant
	ScopeUser
)

// Retencao Aggregate
type Retencao struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Scope     Scope
	Dias      int
	Ativo     bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Repository interface {
	Save(p *Retencao) error
	GetActivePolicy(tenantID uuid.UUID, scope Scope) (*Retencao, error)
}
