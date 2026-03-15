package rabbitmq

import (
	"encoding/json"
	"testing"
	"time"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/google/uuid"
)

func TestPublisher_JobJSONMarshal(t *testing.T) {
	startedAt := time.Now()
	params := `{"resolution":"1080p"}`

	job := &domainJob.JobProcessamento{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		VideoID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Status:    domainJob.StatusQueued,
		Tipo:      1,
		StartedAt: &startedAt,
		Params:    &params,
	}

	body, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("failed to marshal job: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("expected marshalled body to not be empty")
	}

	var unmarshalledJob domainJob.JobProcessamento
	if err := json.Unmarshal(body, &unmarshalledJob); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	if unmarshalledJob.ID != job.ID {
		t.Errorf("expected ID %v, got %v", job.ID, unmarshalledJob.ID)
	}
	if unmarshalledJob.Params == nil || *unmarshalledJob.Params != `{"resolution":"1080p"}` {
		t.Errorf("expected Params %v, got %v", `{"resolution":"1080p"}`, unmarshalledJob.Params)
	}
}

func TestNewRabbitMQConnection_ConnectionRefused(t *testing.T) {
	// Configura o host para uma porta errada para garantir que a conexão falhe
	t.Setenv("RABBITMQ_HOST", "localhost")

	// A conexão deve falhar quando o RabbitMQ não está disponível
	_, _, err := NewRabbitMQConnection()
	if err == nil {
		t.Errorf("expected connection error, got nil")
	} else {
		t.Logf("Got expected connection error: %v", err)
	}
}
