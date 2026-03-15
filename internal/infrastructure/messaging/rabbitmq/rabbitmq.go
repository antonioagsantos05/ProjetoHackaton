package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// NewRabbitMQConnection estabelece uma conexão com o RabbitMQ e retorna a conexão e um canal.
func NewRabbitMQConnection() (*amqp.Connection, *amqp.Channel, error) {
	host := os.Getenv("RABBITMQ_HOST")
	user := "guest"    // Usuário padrão da imagem do RabbitMQ
	password := "guest" // Senha padrão

	connURL := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host)

	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao conectar com o RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // Fecha a conexão se não conseguir abrir o canal
		return nil, nil, fmt.Errorf("erro ao abrir um canal no RabbitMQ: %w", err)
	}

	log.Println("Conexão com o RabbitMQ estabelecida com sucesso!")
	return conn, ch, nil
}
