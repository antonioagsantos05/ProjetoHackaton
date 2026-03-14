package main

import (
	"os"
	"testing"
)

func TestMainEnv_RabbitMQUrl(t *testing.T) {
	// Apenas para simular a var de ambiente
	os.Setenv("RABBITMQ_URL", "amqp://user:pass@localhost:5672/")

	url := os.Getenv("RABBITMQ_URL")
	if url != "amqp://user:pass@localhost:5672/" {
		t.Errorf("expected amqp url from env, got %v", url)
	}
}