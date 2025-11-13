package routers

import (
	"context"
	"go/hioto-logger/config"
	"go/hioto-logger/pkg/handler"
	"go/hioto-logger/pkg/utils"
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
		config.MQTT_LOCAL_INSTANCE_NAME.GetValue(),
		config.MQTT_TOPIC_STATUS_DEVICE.GetValue(),
		r.consumerHandler.ChangeStatusDevice)
}
