package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/fiap-x/video-processor/internal/usecase/video"
	"github.com/google/uuid"
)

type mockVideoRepo struct{}
func (m *mockVideoRepo) Save(v *domainVideo.Video) error { return nil }
func (m *mockVideoRepo) FindByID(id uuid.UUID) (*domainVideo.Video, error) { return nil, nil }
func (m *mockVideoRepo) UpdateStatus(id uuid.UUID, s domainVideo.Status) error { return nil }

type mockJobRepo struct{}
func (m *mockJobRepo) Create(j *domainJob.Job) error { return nil }
func (m *mockJobRepo) Update(j *domainJob.Job) error { return nil }
func (m *mockJobRepo) FindByID(id uuid.UUID) (*domainJob.Job, error) { return nil, nil }

type mockQueue struct{}
func (m *mockQueue) PublishJob(j *domainJob.Job) error { return nil }

func TestVideoHandler_Upload_Success(t *testing.T) {
	uc := video.NewUploadVideoUseCase(&mockVideoRepo{}, &mockJobRepo{}, &mockQueue{})
	handler := NewVideoHandler(uc)

	reqBody := UploadRequest{
		TenantID: uuid.New().String(),
		UserID:   uuid.New().String(),
		Title:    "Test Video",
		Meta:     "{}",
	}
	bodyData, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/videos", bytes.NewBuffer(bodyData))
	rr := httptest.NewRecorder()

	handler.Upload(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	// Criamos um tipo aux para unmarshal sem conflito com as datas
	type VideoResponse struct {
		ID        string `json:"ID"`
		TenantID  string `json:"TenantID"`
		UserID    string `json:"UserID"`
		Title     string `json:"Title"`
		Status    int    `json:"Status"`
		Meta      string `json:"Meta"`
		CreatedAt time.Time `json:"CreatedAt"`
		UpdatedAt time.Time `json:"UpdatedAt"`
	}

	var res VideoResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if res.Title != "Test Video" {
		t.Errorf("expected title %v, got %v", "Test Video", res.Title)
	}
}

func TestVideoHandler_Upload_InvalidJSON(t *testing.T) {
	uc := video.NewUploadVideoUseCase(&mockVideoRepo{}, &mockJobRepo{}, &mockQueue{})
	handler := NewVideoHandler(uc)

	req := httptest.NewRequest("POST", "/videos", bytes.NewBuffer([]byte("{invalid json")))
	rr := httptest.NewRecorder()

	handler.Upload(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %v, got %v", http.StatusBadRequest, rr.Code)
	}
}

func TestVideoHandler_Upload_InvalidUUID(t *testing.T) {
	uc := video.NewUploadVideoUseCase(&mockVideoRepo{}, &mockJobRepo{}, &mockQueue{})
	handler := NewVideoHandler(uc)

	reqBody := UploadRequest{
		TenantID: "invalid-uuid",
		UserID:   uuid.New().String(),
		Title:    "Test Video",
		Meta:     "{}",
	}
	bodyData, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/videos", bytes.NewBuffer(bodyData))
	rr := httptest.NewRecorder()

	handler.Upload(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %v, got %v", http.StatusBadRequest, rr.Code)
	}
}
