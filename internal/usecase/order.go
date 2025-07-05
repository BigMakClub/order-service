package usecase

import (
	"context"
	"github.com/google/uuid"
	"log"
	"order-service/internal/domain"
	"order-service/internal/repository"
)

type OrderUC struct {
	repo  repository.OrderRepository
	cache Cache
}

func NewOrederUC(r repository.OrderRepository, c Cache) *OrderUC {
	return &OrderUC{r, c}
}

func (uc *OrderUC) Get(ctx context.Context, id string) (*domain.Order, error) {
	log.Printf("[OrderUC] Getting order %s", id)
	
	// Try to parse UUID to validate format
	_, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[OrderUC] Invalid UUID format: %s, error: %v", id, err)
		return nil, err
	}

	// Try cache first
	if v, ok := uc.cache.Get(id); ok {
		log.Printf("[OrderUC] Found in cache: %s", id)
		return v, nil
	}

	log.Printf("[OrderUC] Not found in cache, trying DB: %s", id)
	order, err := uc.repo.Find(ctx, id)
	if err != nil {
		log.Printf("[OrderUC] DB error: %v", err)
		return nil, err
	}
	
	if order == nil {
		log.Printf("[OrderUC] Order not found in DB: %s", id)
		return nil, nil
	}

	log.Printf("[OrderUC] Found in DB, saving to cache: %s", id)
	uc.cache.Set(id, order)
	
	return order, nil
}

func (uc *OrderUC) Set(ctx context.Context, order *domain.Order) error {
	orderID := order.OrderId.String()
	log.Printf("[OrderUC] Saving order %s", orderID)
	
	if err := uc.repo.Save(ctx, order); err != nil {
		log.Printf("[OrderUC] Error saving to DB: %v", err)
		return err
	}
	
	log.Printf("[OrderUC] Successfully saved to DB: %s", orderID)
	uc.cache.Set(orderID, order)
	log.Printf("[OrderUC] Successfully saved to cache: %s", orderID)
	return nil
}
