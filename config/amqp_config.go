package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQInstance struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

var rmqInstances = make(map[string]*RMQInstance)
var mu sync.Mutex

func initializeRabbitMQ(url, rmqInstance string) error {
	mu.Lock()
	defer mu.Unlock()

	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := range 5 {
		conn, err = amqp.Dial(url)

		if err == nil {
			log.Printf("✅ Successfully connected to RabbitMQ (%s)", rmqInstance)

			ch, err = conn.Channel()
			if err != nil {
				log.Printf("❌ Failed to open channel: %v", err)
				return err
			}

			rmqInstances[rmqInstance] = &RMQInstance{
				Conn:    conn,
				Channel: ch,
			}

			log.Println("✅ RabbitMQ channel opened successfully")
			return nil
		}

		log.Printf("⚠️ Failed to connect to RabbitMQ, retrying in 5 seconds... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	return err
}

func GetRMQInstance(rmqtype string) (*RMQInstance, error) {
	mu.Lock()
	defer mu.Unlock()

	instance, ok := rmqInstances[rmqtype]

	if !ok {
		return nil, fmt.Errorf("RabbitMQ instance %s not found", rmqtype)
	}

	return instance, nil
}

func CloseRabbitMQ() {
	mu.Lock()
	defer mu.Unlock()

	for name, rmqInstance := range rmqInstances {
		if rmqInstance.Channel != nil {
			rmqInstance.Channel.Close()
			log.Printf("🔒 RabbitMQ %s channel closed \n", name)
		}

		if rmqInstance.Conn != nil {
			rmqInstance.Conn.Close()
			log.Printf("🔒 RabbitMQ %s connection closed", name)
		}
	}
}

func CreateRmqInstance() {
	if err := initializeRabbitMQ(os.Getenv("RMQ_HIOTO_CLOUD_URL"), os.Getenv("RMQ_HIOTO_CLOUD_INSTANCE")); err != nil {
		log.Fatal(err)
	}
}
