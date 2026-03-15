package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/fiap-x/video-processor/internal/usecase/video"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Mocks

type mockUploadUseCase struct{}

func (m *mockUploadUseCase) Execute(filePath string, userID, tenantID uuid.UUID) (*domainJob.JobProcessamento, error) {
	return &domainJob.JobProcessamento{ID: uuid.New()}, nil
}

type mockListUseCase struct {
	videos []video.VideoDTO
}

func (m *mockListUseCase) Execute(userID string) ([]video.VideoDTO, error) {
	if m.videos != nil {
		return m.videos, nil
	}
	return []video.VideoDTO{}, nil
}

func TestVideoHandler_ListVideos_Success(t *testing.T) {
	mockList := &mockListUseCase{
		videos: []video.VideoDTO{
			{ID: uuid.New(), Title: "Video 1", Status: 0},
			{ID: uuid.New(), Title: "Video 2", Status: 2},
		},
	}
	handler := NewVideoHandler(&mockUploadUseCase{}, mockList)

	req := httptest.NewRequest("GET", "/users/"+uuid.New().String()+"/videos", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/videos", handler.ListVideos).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res []video.VideoDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("expected 2 videos, got %v", len(res))
	}
}

func TestVideoHandler_ListVideos_Empty(t *testing.T) {
	handler := NewVideoHandler(&mockUploadUseCase{}, &mockListUseCase{})

	req := httptest.NewRequest("GET", "/users/"+uuid.New().String()+"/videos", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/videos", handler.ListVideos).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
}
