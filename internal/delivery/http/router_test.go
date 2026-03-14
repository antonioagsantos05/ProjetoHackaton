package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter_HealthCheck(t *testing.T) {
	router := NewRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	// Como tem middleware de auth, o health route precisa de um token se estiver ativado globalmente.
	req.Header.Set("Authorization", "Bearer valid_token")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	expected := `{"status": "UP"}`
	if rr.Body.String() != expected {
		t.Errorf("expected body %v, got %v", expected, rr.Body.String())
	}
}

func TestNewRouter_MiddlewareAuthApplied(t *testing.T) {
	router := NewRouter()

	// Tentativa de acesso sem token
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, rr.Code)
	}
}
