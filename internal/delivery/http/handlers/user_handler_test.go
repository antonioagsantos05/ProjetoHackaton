package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domainUser "github.com/fiap-x/video-processor/internal/domain/user"
	"github.com/fiap-x/video-processor/internal/usecase/user"
	"github.com/google/uuid"
)

// Mocks

type mockCreateUserUseCase struct{}

func (m *mockCreateUserUseCase) Execute(req user.CreateUserRequest) (*domainUser.User, error) {
	return &domainUser.User{
		ID:       uuid.New(),
		TenantID: uuid.New(),
		Email:    req.Email,
		Nome:     req.Nome,
		Status:   1,
	}, nil
}

type mockLoginUseCase struct{}

func (m *mockLoginUseCase) Execute(email, password string) (string, error) {
	return "mock-jwt-token", nil
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	handler := NewUserHandler(&mockCreateUserUseCase{}, &mockLoginUseCase{})

	body, _ := json.Marshal(user.CreateUserRequest{
		Email:    "teste@exemplo.com",
		Nome:     "Teste",
		Password: "senha123",
	})

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	var res domainUser.User
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res.Email != "teste@exemplo.com" {
		t.Errorf("expected email teste@exemplo.com, got %v", res.Email)
	}
}

func TestUserHandler_Login_Success(t *testing.T) {
	handler := NewUserHandler(&mockCreateUserUseCase{}, &mockLoginUseCase{})

	body, _ := json.Marshal(LoginRequest{
		Email:    "teste@exemplo.com",
		Password: "senha123",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["token"] == "" {
		t.Errorf("expected token, got empty")
	}
}
