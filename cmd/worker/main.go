package main

import (
	"log"

	"github.com/fiap-x/video-processor/internal/worker"
	"github.com/streadway/amqp"
)

func main() {
	// A URI real viria de variáveis de ambiente
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Iniciando Worker de Processamento de Vídeos FIAP X...")
	worker.StartVideoProcessor(conn)
}
