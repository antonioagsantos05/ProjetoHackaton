package worker

import (
	"log"

	"github.com/streadway/amqp"
)

// StartVideoProcessor inicia o listener do RabbitMQ simulando um Worker de FFmpeg
func StartVideoProcessor(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir o canal do rabbitmq: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"video_processing_jobs", // Fila
		"",                      // Consumer
		false,                   // Auto-Ack (Manual é melhor para retries e DLQ)
		false,                   // Exclusive
		false,                   // No-local
		false,                   // No-Wait
		nil,                     // Args
	)
	if err != nil {
		log.Fatalf("Falha ao registrar consumidor: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Worker recebeu Job: %s", d.Body)

			// TODO: Deserializar o Job JSON
			// TODO: Baixar o Vídeo do S3
			// TODO: Executar FFmpeg para extrair frames via comando exec.Command
			// TODO: Gerar ZIP
			// TODO: Atualizar Job e Video Status no Banco
			// TODO: Enviar notificação

			log.Printf("Processamento finalizado. Enviando ACK.")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Worker inicializado. Aguardando mensagens. Para sair pressione CTRL+C")
	<-forever
}
