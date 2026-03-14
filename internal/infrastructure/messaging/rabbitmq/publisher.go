package rabbitmq

import (
	"encoding/json"

	domainJob "github.com/fiap-x/video-processor/internal/domain/job"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declarar a fila (exchange, queue) aqui em um cenário real
	_, err = ch.QueueDeclare(
		"video_processing_jobs",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *Publisher) PublishJob(job *domainJob.Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",                      // exchange
		"video_processing_jobs", // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
}
