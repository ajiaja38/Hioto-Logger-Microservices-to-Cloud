package main

import (
	"fmt"
	"go/hioto-logger/config"
	"go/hioto-logger/pkg/service"
	"log"
	"os"
	"os/signal"
	"sync"
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

	if err := config.InitializeRabbitMQ(os.Getenv("RMQ_HIOTO_CLOUD_URL"), "Hioto Biznet"); err != nil {
		log.Fatal(err)
	}

	defer config.CloseRabbitMQ()

	logService := service.NewLogService(db)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var mu sync.Mutex
	running := false

	log.Println("🚀 Logger started, publishing every 10 seconds...")

	for {
		select {
		case <-stop:
			log.Println("🛑 Shutting down gracefully...")
			return

		case t := <-ticker.C:
			mu.Lock()
			if running {
				log.Println("⏳ Previous job still running, skip.")
				mu.Unlock()
				continue
			}
			running = true
			mu.Unlock()

			go func() {
				defer func() {
					mu.Lock()
					running = false
					mu.Unlock()
				}()

				fmt.Printf("\n⏰ Tick at: %s\n", t.Format("15:04:05"))
				logService.GetAllLogs()
				logService.GetAllLogAktuators()
				logService.GetAllMonitoringHistory()
			}()
		}
	}
}
