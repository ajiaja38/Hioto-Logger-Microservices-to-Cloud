package main

import (
	"context"
	"go/hioto-logger/config"
	"go/hioto-logger/pkg/handler"
	"go/hioto-logger/pkg/routers"
	"go/hioto-logger/pkg/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Print("Error loading .env file")
	}
}

func main() {
	db, err := config.DBConnection()

	if err != nil {
		log.Print(err)
	}

	config.CreateRmqInstance()
	defer config.CloseRabbitMQ()

	// service
	logService := service.NewLogService(db)
	deviceService := service.NewDeviceService(db)

	// handler
	consumerHandler := handler.NewConsumerHandler(deviceService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// router
	consumerRouter := routers.NewConsumerMessageBroker(consumerHandler, ctx)
	consumerRouter.StartConsume()

	service.CronJobService(
		time.NewTicker(time.Second*10),
		logService.GetAllLogs,
		logService.GetAllLogAktuators,
		logService.GetAllMonitoringHistory,
	)
}
