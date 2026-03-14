package rabbitmq

import (
	"encoding/json"
	"testing"
	"time"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Devido a dependência externa com o cluster RabbitMQ, testes unitários clássicos para drivers
// externos de AMQP envolvem Mocks de Connection/Channel, que no caso do streadway/amqp exigem
// wrappers customizados. Como alternativa faremos um teste que verifica a serialização do payload do JSON
// no Job, que é o que o Publisher processa localmente antes de enviar para o Channel.

func TestPublisher_JobJSONMarshal(t *testing.T) {
	startedAt := time.Now()

	job := &domainJob.Job{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		TenantID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		VideoID:    uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Status:     domainJob.StatusQueued,
		Type:       1,
		StartedAt:  &startedAt,
		Params:     `{"resolution":"1080p"}`,
	}

	body, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("failed to marshal job: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("expected marshalled body to not be empty")
	}

	var unmarshalledJob domainJob.Job
	if err := json.Unmarshal(body, &unmarshalledJob); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	if unmarshalledJob.ID != job.ID {
		t.Errorf("expected ID %v, got %v", job.ID, unmarshalledJob.ID)
	}
	if unmarshalledJob.Params != `{"resolution":"1080p"}` {
		t.Errorf("expected Params %v, got %v", `{"resolution":"1080p"}`, unmarshalledJob.Params)
	}
}

// Este teste irá falhar rápido de forma esperada quando não há um RabbitMQ mockado ou subindo na URL
func TestNewPublisher_ConnectionRefused(t *testing.T) {
	_, err := NewPublisher("amqp://guest:guest@localhost:9999/") // Porta errada / Mock

	if err == nil {
		t.Errorf("expected connection error, got nil")
	} else if err != amqp.ErrClosed && err.Error() != "dial tcp [::1]:9999: connectex: No connection could be made because the target machine actively refused it." && err.Error() != "dial tcp 127.0.0.1:9999: connect: connection refused" {
		// Log da falha real dependendo do SO/network do test runner, mas a gente só quer garantir que dá erro
		t.Logf("Got expected connection error: %v", err)
	}
}
