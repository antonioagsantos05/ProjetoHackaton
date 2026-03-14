package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestDownloadHandler_GetZip_Success(t *testing.T) {
	handler := NewDownloadHandler()

	req := httptest.NewRequest("GET", "/downloads/abc-123", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/downloads/{id}", handler.GetZip).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedUrl := "https://s3.aws.com/fiapx/abc-123.zip"
	if res["url"] != expectedUrl {
		t.Errorf("expected url %v, got %v", expectedUrl, res["url"])
	}
}
