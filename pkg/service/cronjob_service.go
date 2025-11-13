package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func CronJobService(timeTicker *time.Ticker, callbacks ...func()) {
	ticker := timeTicker
	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var mu sync.Mutex
	running := false

	log.Println("🚀 Cronjob started...")

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
					if r := recover(); r != nil {
						log.Printf("🔥 Panic recovered in cron job: %v", r)
					}
					mu.Lock()
					running = false
					mu.Unlock()
				}()

				fmt.Printf("\n⏰ Tick at: %s\n", t.Format("15:04:05"))
				for _, callback := range callbacks {
					callback()
				}
			}()

		}
	}
}
