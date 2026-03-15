package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fiap-x/video-processor/internal/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserRequest define os dados necessários para criar um usuário.
type CreateUserRequest struct {
	Email    string `json:"email"`
	Nome     string `json:"nome"`
	Password string `json:"password"`
}

// CreateUserUseCase define a interface para o caso de uso de criação de usuário.
type CreateUserUseCase interface {
	Execute(req CreateUserRequest) (*user.User, error)
}

// createUserUseCase é a implementação de CreateUserUseCase.
type createUserUseCase struct {
	db *sql.DB
}

// NewCreateUserUseCase cria uma nova instância de createUserUseCase.
func NewCreateUserUseCase(db *sql.DB) CreateUserUseCase {
	return &createUserUseCase{db: db}
}

// Execute orquestra a criação de um novo usuário.
func (uc *createUserUseCase) Execute(req CreateUserRequest) (*user.User, error) {
	// 1. Criptografar a senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao criptografar a senha: %w", err)
	}

	// 2. Criar a entidade do usuário
	newUser := &user.User{
		ID:        uuid.New(),
		TenantID:  uuid.New(), // Simplificação: usando um novo tenant para cada usuário
		Email:     req.Email,
		Nome:      req.Nome,
		HashSenha: string(hashedPassword),
		Status:    1, // Ativo
		CreatedAt: time.Now(),
	}

	// 3. Inserir no banco de dados
	query := `INSERT INTO usuario (id, tenant_id, email, nome, hash_senha, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = uc.db.Exec(query, newUser.ID, newUser.TenantID, newUser.Email, newUser.Nome, newUser.HashSenha, newUser.Status, newUser.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir usuário no banco de dados: %w", err)
	}

	return newUser, nil
}
