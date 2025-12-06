package utils

import (
	"go/hioto-logger/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishToRmq(rmqInstance string, message []byte, queueName string, exchange string) {
	instance, err := config.GetRMQInstance(config.RMQ_CLOUD_INSTANCE.GetValue())

	if err != nil {
		log.Println(err)
		return
	}

	ch := instance.Channel

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
		log.Printf("Failed to declare the queue: %v 💥", err)
		return
	}

	if err = ch.Publish(
		exchange,
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body: message,
		},
	); err != nil {
		log.Printf("Failed to publish data: %v 💥", err)
		return
	}

	log.Printf("Published message to queue %s ✅", queueName)
}
