package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createTestToken() string {
	claims := jwt.MapClaims{
		"sub": "00000000-0000-0000-0000-000000000001",
		"tid": "00000000-0000-0000-0000-000000000002",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(""))
	return tokenString
}

func TestNewRouter_HealthCheck(t *testing.T) {
	router := NewRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	req.Header.Set("Authorization", "Bearer "+createTestToken())

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

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, rr.Code)
	}
}
