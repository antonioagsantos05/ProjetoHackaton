package user

import (
	"time"

	"github.com/google/uuid"
)

// User Aggregate Root
type User struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Email     string
	Name      string
	HashPass  string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Repository define a interface de persistência para User
type Repository interface {
	Save(user *User) error
	FindByEmail(email string) (*User, error)
}
