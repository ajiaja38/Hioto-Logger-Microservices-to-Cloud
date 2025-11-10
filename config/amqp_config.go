package config

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RmqConn *amqp.Connection
var RmqChannel *amqp.Channel

func InitializeRabbitMQ(url, rmqType string) error {
	var err error

	for i := 0; i < 5; i++ {
		RmqConn, err = amqp.Dial(url)
		if err == nil {
			log.Printf("✅ Successfully connected to RabbitMQ (%s)", rmqType)

			RmqChannel, err = RmqConn.Channel()
			if err != nil {
				log.Printf("❌ Failed to open channel: %v", err)
				return err
			}

			log.Println("✅ RabbitMQ channel opened successfully")
			return nil
		}

		log.Printf("⚠️ Failed to connect to RabbitMQ, retrying in 5 seconds... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	return err
}

func CloseRabbitMQ() {
	if RmqChannel != nil {
		RmqChannel.Close()
		log.Println("🔒 RabbitMQ channel closed")
	}

	if RmqConn != nil {
		RmqConn.Close()
		log.Println("🔒 RabbitMQ connection closed")
	}
}
