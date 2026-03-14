package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fiap-x/video-processor/internal/usecase/user"
)

type UserHandler struct {
	createUserUseCase user.CreateUserUseCase
	loginUseCase      user.LoginUseCase
}

func NewUserHandler(createUC user.CreateUserUseCase, loginUC user.LoginUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUC,
		loginUseCase:      loginUC,
	}
}

// LoginRequest define o corpo da requisição para o endpoint de login.
type LoginRequest struct {
	Email    string `json:"email" example:"teste@exemplo.com"`
	Password string `json:"password" example:"senha_forte_123"`
}

// CreateUser godoc
// @Summary      Cria um novo usuário
// @Description  Registra um novo usuário no sistema.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      user.CreateUserRequest  true  "Informações do Usuário"
// @Success      201   {object}  user.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.createUserUseCase.Execute(req)
	if err != nil {
		log.Printf("Erro ao executar o use case de criação de usuário: %v", err)
		http.Error(w, "Erro ao criar usuário.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// Login godoc
// @Summary      Autentica um usuário
// @Description  Realiza o login de um usuário e retorna um token JWT.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Credenciais de Login"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Chamar o Use Case de Login
	token, err := h.loginUseCase.Execute(req.Email, req.Password)
	if err != nil {
		log.Printf("Falha na tentativa de login para o email '%s': %v", req.Email, err)
		http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	// Se o login for bem-sucedido, retorna o token real
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
