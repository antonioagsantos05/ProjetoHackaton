package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestStatusHandler_CheckJobStatus_Success(t *testing.T) {
	handler := NewStatusHandler()

	req := httptest.NewRequest("GET", "/jobs/123/status", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/jobs/{id}/status", handler.CheckJobStatus).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["id"] != "123" {
		t.Errorf("expected id %v, got %v", "123", res["id"])
	}
	if res["status"] != "PROCESSING" {
		t.Errorf("expected status %v, got %v", "PROCESSING", res["status"])
	}
}
