package kafkaC

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"order-service/internal/domain"
	"order-service/internal/usecase"
	"time"
)

func Start(
	ctx context.Context,
	uc *usecase.OrderUC,
	brokers []string,
	topic string,
	groupID string,
) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MinBytes:       1_000,
		MaxBytes:       1_000_000,
		MaxWait:        500 * time.Millisecond,
	})

	errCh := make(chan error, 1)
	go func() {
		defer r.Close()
		defer close(errCh)

		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("kafka read: %v", err)
				errCh <- fmt.Errorf("kafka read error: %w", err)
				return
			}

			var ord domain.Order
			if err := json.Unmarshal(msg.Value, &ord); err != nil {
				log.Printf("bad JSON, skip: %v", err)
				_ = r.CommitMessages(ctx, msg) // ack «мусор»
				continue
			}

			if err := uc.Set(context.Background(), &ord); err != nil {
				fmt.Printf("save error: %v\n", err)
				log.Printf("save failed (retry later): %v", err)
				continue
			}

			if err := r.CommitMessages(ctx, msg); err != nil {
				log.Printf("commit offset: %v", err)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
