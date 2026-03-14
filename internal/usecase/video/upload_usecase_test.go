package video

import (
	"errors"
	"testing"
	"time"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/google/uuid"
)

// Mocks

type mockVideoRepo struct {
	saveFunc func(video *domainVideo.Video) error
}

func (m *mockVideoRepo) Save(video *domainVideo.Video) error {
	if m.saveFunc != nil {
		return m.saveFunc(video)
	}
	return nil
}
func (m *mockVideoRepo) FindByID(id uuid.UUID) (*domainVideo.Video, error)     { return nil, nil }
func (m *mockVideoRepo) UpdateStatus(id uuid.UUID, status domainVideo.Status) error { return nil }

type mockJobRepo struct {
	createFunc func(job *domainJob.Job) error
}

func (m *mockJobRepo) Create(job *domainJob.Job) error {
	if m.createFunc != nil {
		return m.createFunc(job)
	}
	return nil
}
func (m *mockJobRepo) Update(job *domainJob.Job) error          { return nil }
func (m *mockJobRepo) FindByID(id uuid.UUID) (*domainJob.Job, error) { return nil, nil }

type mockQueue struct {
	publishFunc func(job *domainJob.Job) error
}

func (m *mockQueue) PublishJob(job *domainJob.Job) error {
	if m.publishFunc != nil {
		return m.publishFunc(job)
	}
	return nil
}

func TestUploadVideoUseCase_Execute_Success(t *testing.T) {
	vr := &mockVideoRepo{}
	jr := &mockJobRepo{}
	q := &mockQueue{}

	uc := NewUploadVideoUseCase(vr, jr, q)

	tenantID := uuid.New()
	userID := uuid.New()
	title := "Meu Video Teste"
	meta := `{"resolution":"1080p"}`

	video, err := uc.Execute(tenantID, userID, title, meta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if video == nil {
		t.Fatalf("expected video, got nil")
	}
	if video.Title != title {
		t.Errorf("expected title %v, got %v", title, video.Title)
	}
	if video.Status != domainVideo.StatusPending {
		t.Errorf("expected status %v, got %v", domainVideo.StatusPending, video.Status)
	}
	if video.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set")
	}
	if video.UpdatedAt.IsZero() {
		t.Errorf("expected UpdatedAt to be set")
	}
}

func TestUploadVideoUseCase_Execute_VideoRepoError(t *testing.T) {
	expectedErr := errors.New("db error")
	vr := &mockVideoRepo{
		saveFunc: func(video *domainVideo.Video) error {
			return expectedErr
		},
	}
	jr := &mockJobRepo{}
	q := &mockQueue{}

	uc := NewUploadVideoUseCase(vr, jr, q)

	_, err := uc.Execute(uuid.New(), uuid.New(), "Title", "{}")

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestUploadVideoUseCase_Execute_JobRepoError(t *testing.T) {
	expectedErr := errors.New("job db error")
	vr := &mockVideoRepo{}
	jr := &mockJobRepo{
		createFunc: func(job *domainJob.Job) error {
			if job.CreatedAt.IsZero() {
				t.Errorf("expected job CreatedAt to be set")
			}
			return expectedErr
		},
	}
	q := &mockQueue{}

	uc := NewUploadVideoUseCase(vr, jr, q)

	_, err := uc.Execute(uuid.New(), uuid.New(), "Title", "{}")

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestUploadVideoUseCase_Execute_QueueError(t *testing.T) {
	expectedErr := errors.New("queue error")
	vr := &mockVideoRepo{}
	jr := &mockJobRepo{}
	q := &mockQueue{
		publishFunc: func(job *domainJob.Job) error {
			return expectedErr
		},
	}

	uc := NewUploadVideoUseCase(vr, jr, q)

	_, err := uc.Execute(uuid.New(), uuid.New(), "Title", "{}")

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
