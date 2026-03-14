package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserHandler_CreateUser_Success(t *testing.T) {
	handler := NewUserHandler()

	req := httptest.NewRequest("POST", "/users", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	var res map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["message"] != "User created successfully" {
		t.Errorf("expected message 'User created successfully', got %v", res["message"])
	}

	if res["id"] == "" {
		t.Errorf("expected id to be not empty")
	}
}
