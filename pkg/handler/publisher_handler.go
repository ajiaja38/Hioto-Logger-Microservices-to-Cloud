package handler

import (
	"go/hioto-logger/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishToRmq(message []byte, queueName string, exchange string) {
	ch := config.RmqChannel

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": int32(120000),
		},
	)

	if err != nil {
		log.Printf("Failed to declare queue: %v 💥", err)
		return
	}

	err = ch.Publish(
		exchange,
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body: message,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v 💥", err)
		return
	}

	log.Printf("Published message to queue %s ✅", queueName)
}
