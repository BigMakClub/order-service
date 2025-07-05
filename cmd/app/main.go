package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"order-service/cmd/http"
	"order-service/cmd/kafkaC"
	"order-service/internal/adapter/cache"
	"order-service/internal/adapter/db"
	"order-service/internal/usecase"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  .env file not found or failed to load, fallback to OS env")
	}

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	pgDsn := os.Getenv("PG_DSN")
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8081"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Run database migrations
	if err := db.RunMigrations(pgDsn); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	repo, err := db.NewPgRepo(pgDsn)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	mem := cache.NewMemCache(100)
	uc := usecase.NewOrederUC(repo, mem)

	if list, err := repo.CacheRestore(ctx); err == nil {
		for _, o := range list {
			mem.Set(o.OrderId.String(), o)
		}
		log.Printf("cache warmed: %d orders", len(list))
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		if err := http.Start(ctx, uc, httpAddr); err != nil {
			log.Printf("HTTP server error: %v", err)
			stop() // Signal shutdown on error
		}
	}()

	// Start Kafka consumer
	go func() {
		defer wg.Done()
		if err := kafkaC.Start(ctx, uc, brokers, topic, groupID); err != nil {
			log.Printf("Kafka consumer error: %v", err)
			stop() // Signal shutdown on error
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down...")
	
	// Wait for servers to close
	wg.Wait()
	log.Println("Graceful shutdown complete")
}
