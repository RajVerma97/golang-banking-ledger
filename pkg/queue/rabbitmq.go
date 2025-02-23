package queue

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func InitRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	if rabbitMQURI == "" {
		log.Fatal("RABBITMQ_URI is not set in environment variables")
	}

	conn, err := amqp.Dial(rabbitMQURI)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		"transaction_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	log.Println("RabbitMQ connected and queue declared")

	return conn, ch, nil
}
