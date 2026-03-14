package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
)

type UserHandler struct {}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// CreateUser godoc
// @Summary      Cria um novo usuário
// @Description  Registra um novo usuário no sistema.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      object  true  "Informações do Usuário"
// @Success      201   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Auth UseCase com BCrypt e JWT (Requisito: O Sistema deve ser protegido por usuário e senha)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": uuid.New().String(),
		"message": "User created successfully",
	})
}

// Login godoc
// @Summary      Autentica um usuário
// @Description  Realiza o login de um usuário e retorna um token JWT.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body      object  true  "Credenciais de Login"
// @Success      200          {object}  map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"message": "Login successful",
	})
}
