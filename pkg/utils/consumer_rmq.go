package utils

import (
	"context"
	"go/hioto-logger/config"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler func([]byte)

func ConsumeRmq(ctx context.Context, instanceName, queueName string, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Consumer stopped before connection", queueName)
			return
		default:
		}

		instance, err := config.GetRMQInstance(instanceName)

		if err != nil {
			log.Println(err)
			return
		}

		ch := instance.Channel

		if instance == nil || instance.Channel == nil {
			log.Printf("[%s] Channel not ready, retrying...", queueName)
			time.Sleep(5 * time.Second)
			continue
		}

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
			log.Printf("[%s] Queue declare error: %v", queueName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Printf("[%s] Failed to consume: %v", queueName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("[%s] Waiting for messages on Queue [%s]⚡️", instanceName, queueName)

		jobs := make(chan []byte, 100)
		wg := &sync.WaitGroup{}
		for range 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for body := range jobs {
					handler(body)
				}
			}()
		}

	consumeLoop:
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] Stopping consumer...", queueName)
				break consumeLoop
			case d, ok := <-msgs:
				if !ok {
					log.Printf("[%s] Message channel closed", queueName)
					break consumeLoop
				}
				select {
				case jobs <- d.Body:
				case <-ctx.Done():
					break consumeLoop
				}
			}
		}

		close(jobs)
		wg.Wait()

		if ctx.Err() != nil {
			return
		}

		log.Printf("[%s] Reconnecting after disconnect...", queueName)
		time.Sleep(5 * time.Second)
	}
}
