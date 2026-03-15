package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Mock

type mockDownloadUseCase struct{}

func (m *mockDownloadUseCase) Execute(jobID string) (string, error) {
	return "https://minio:9000/videos/abc-123.zip?token=presigned", nil
}

func TestDownloadHandler_GetZip_Success(t *testing.T) {
	handler := NewDownloadHandler(&mockDownloadUseCase{})

	req := httptest.NewRequest("GET", "/downloads/abc-123", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/downloads/{id}", handler.GetZip).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res["url"] == "" {
		t.Errorf("expected url to be present")
	}
}
