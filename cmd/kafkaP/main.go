package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"log"
	"math/rand"
	"order-service/internal/domain"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found or failed to load, fallback to OS env")
	}
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		fmt.Println("KAFKA_BROKERS environment variable not set")
		brokers = "localhost:29092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "orders-topic"
		fmt.Println("KAFKA_BROKERS environment variable not set")
	}
	interval := time.Second

	// ── writer ──────────────────────────────────────────────────
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}
	defer w.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Printf("emitter → %s (%s) every %s", brokers, topic, interval)

	for n := 0; ; n++ {
		select {
		case <-ctx.Done():
			log.Println("stopped")
			return
		default:
			ord := makeDummyOrder(n)
			payload, _ := json.Marshal(ord)

			msg := kafka.Message{
				Key:   []byte(ord.OrderId.String()),
				Value: payload,
				Time:  time.Now(),
			}
			if err := w.WriteMessages(ctx, msg); err != nil {
				log.Printf("write: %v", err)
			} else {
				log.Printf("sent %s", ord.OrderId)
			}
			time.Sleep(interval)
		}
	}
}

func makeDummyOrder(n int) domain.Order {
	uid := uuid.New()
	return domain.Order{
		OrderId:     uid,
		TrackNumber: fmt.Sprintf("TESTTRACK-%06d", n),
		Entry:       "WBIL",
		Locale:      "en",
		CustomerId:  "emitter",
		ShardKey:    int64(rand.Intn(1000)),
		SmId:        rand.Intn(1000),
		DateCreated: time.Now(),
		Delivery: domain.Delivery{
			Name:    "Load Gen",
			Phone:   "+100000000",
			City:    "GoCity",
			Address: "Benchmark str.",
		},
		Payment: domain.Payment{
			TransactionId: uid.String(),
			Amount:        int64(rand.Intn(5000) + 500),
			Currency:      "USD",
			Provider:      "emitter-pay",
		},
		Items: []domain.Item{{
			ChrtId:      int64(rand.Intn(1e7)),
			TrackNumber: "TESTTRACK",
			Name:        "DemoItem",
			Price:       999,
			TotalPrice:  999,
			Brand:       "EmitterCo",
			Status:      "202",
		}},
	}
}
