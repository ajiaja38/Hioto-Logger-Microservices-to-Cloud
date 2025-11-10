package config

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RmqConnection(uri string, rmqType string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := range 5 {
		conn, err = amqp.Dial(uri)

		if err == nil {
			log.Printf("successfully connected to rabbitMQ %s🔀", rmqType)
			return conn, nil
		}

		log.Printf("failed to connect to rabbitMQ, retrying in 5 seconds... (%d/5) 💥", i+1)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}
