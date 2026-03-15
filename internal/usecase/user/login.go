package user

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fiap-x/video-processor/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Definindo uma chave secreta para assinar o token.
// Em produção, isso DEVE vir de uma variável de ambiente segura.
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// LoginUseCase define a interface para o caso de uso de login.
type LoginUseCase interface {
	Execute(email, password string) (string, error) // Retorna o token em vez do usuário
}

// loginUseCase é a implementação de LoginUseCase.
type loginUseCase struct {
	db *sql.DB
}

// NewLoginUseCase cria uma nova instância de loginUseCase.
func NewLoginUseCase(db *sql.DB) LoginUseCase {
	return &loginUseCase{db: db}
}

// Execute orquestra a autenticação de um usuário e retorna um token JWT.
func (uc *loginUseCase) Execute(email, password string) (string, error) {
	// 1. Buscar usuário pelo e-mail
	var u user.User
	query := `SELECT id, tenant_id, email, nome, hash_senha, status FROM usuario WHERE email = $1`
	err := uc.db.QueryRow(query, email).Scan(&u.ID, &u.TenantID, &u.Email, &u.Nome, &u.HashSenha, &u.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("usuário ou senha inválidos")
		}
		return "", fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	// 2. Comparar a senha
	err = bcrypt.CompareHashAndPassword([]byte(u.HashSenha), []byte(password))
	if err != nil {
		return "", fmt.Errorf("usuário ou senha inválidos")
	}

	// 3. Verificar se o usuário está ativo
	if u.Status != 1 {
		return "", fmt.Errorf("usuário inativo")
	}

	// 4. Gerar o token JWT
	claims := jwt.MapClaims{
		"sub": u.ID,         // "Subject", o ID do usuário
		"tid": u.TenantID,   // "Tenant ID"
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expira em 24 horas
		"iat": time.Now().Unix(),                     // "Issued At"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar o token JWT: %w", err)
	}

	return tokenString, nil
}
