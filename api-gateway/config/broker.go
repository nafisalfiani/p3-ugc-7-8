package config

import (
	"log"

	"github.com/streadway/amqp"
)

func InitBroker() (*amqp.Connection, error) {
	// RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	return conn, nil
}
