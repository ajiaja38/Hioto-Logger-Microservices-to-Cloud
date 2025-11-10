package routers

import (
	"context"
	"go/hioto-logger/pkg/handler"
	"go/hioto-logger/pkg/utils"
	"os"
)

type ConsumerMessageBroker struct {
	consumerHandler *handler.ConsumerHandler
	ctx             context.Context
}

func NewConsumerMessageBroker(consumerHandler *handler.ConsumerHandler, ctx context.Context) *ConsumerMessageBroker {
	return &ConsumerMessageBroker{
		consumerHandler: consumerHandler,
		ctx:             ctx,
	}
}

func (r *ConsumerMessageBroker) StartConsume() {
	go utils.ConsumeMQTTTopic(
		r.ctx,
		os.Getenv("MQTT_LOCAL_INSTANCE_NAME"),
		os.Getenv("STATUS_DEVICE_TOPIC"),
		r.consumerHandler.ChangeStatusDevice)
}
