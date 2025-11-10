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
	go utils.ConsumeRmq(r.ctx, os.Getenv("RMQ_HIOTO_LOCAL_INSTANCE"), os.Getenv("STATUS_DEVICE_QUEUE"), r.consumerHandler.ChangeStatusDevice)
}
