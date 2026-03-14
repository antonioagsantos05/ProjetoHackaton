package video

import (
	"time"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	domainVideo "github.com/fiap-x/video-processor/internal/domain/video"
	"github.com/google/uuid"
)

type UploadVideoUseCase struct {
	videoRepo domainVideo.Repository
	jobRepo   domainJob.Repository
	queue     domainJob.MessageQueue
}

func NewUploadVideoUseCase(
	vr domainVideo.Repository,
	jr domainJob.Repository,
	q domainJob.MessageQueue,
) *UploadVideoUseCase {
	return &UploadVideoUseCase{
		videoRepo: vr,
		jobRepo:   jr,
		queue:     q,
	}
}

// Execute orquestra o upload do vídeo, persistência e envio para a fila
func (uc *UploadVideoUseCase) Execute(tenantID, userID uuid.UUID, title string, meta string) (*domainVideo.Video, error) {
	now := time.Now()

	// 1. Criar entidade de Vídeo
	video := &domainVideo.Video{
		ID:        uuid.New(),
		TenantID:  tenantID,
		UserID:    userID,
		Title:     title,
		Status:    domainVideo.StatusPending,
		Meta:      meta,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 2. Persistir Vídeo no banco (PostgreSQL)
	if err := uc.videoRepo.Save(video); err != nil {
		return nil, err
	}

	// 3. Criar Job de Processamento
	job := &domainJob.Job{
		ID:        uuid.New(),
		TenantID:  tenantID,
		VideoID:   video.ID,
		Type:      1, // 1 = Extração de Frames
		Status:    domainJob.StatusQueued,
		Params:    "{}", // Parametros default (pode vir do meta no futuro)
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 4. Salvar Job
	if err := uc.jobRepo.Create(job); err != nil {
		return nil, err
	}

	// 5. Publicar evento/job na fila (RabbitMQ/Kafka) para os Workers
	if err := uc.queue.PublishJob(job); err != nil {
		return nil, err
	}

	return video, nil
}
