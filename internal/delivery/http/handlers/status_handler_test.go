package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/fiap-x/video-processor/internal/usecase/job"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Mock

type mockStatusUseCase struct{}

func (m *mockStatusUseCase) Execute(jobID string) (*job.StatusResponse, error) {
	return &job.StatusResponse{
		JobInfo: domainJob.JobProcessamento{
			ID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Status: domainJob.StatusProcessing,
		},
		VideoInfo: domainVideo.Video{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Title: "test.mp4",
		},
	}, nil
}

func TestStatusHandler_CheckJobStatus_Success(t *testing.T) {
	handler := NewStatusHandler(&mockStatusUseCase{})

	req := httptest.NewRequest("GET", "/jobs/00000000-0000-0000-0000-000000000001/status", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/jobs/{id}/status", handler.CheckJobStatus).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var res job.StatusResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res.JobInfo.Status != domainJob.StatusProcessing {
		t.Errorf("expected job status PROCESSING, got %v", res.JobInfo.Status)
	}
}
