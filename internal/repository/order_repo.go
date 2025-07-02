package repository

import (
	"context"
	"order-service/order-service/internal/domain"
)

type OrderRepository interface {
	Find(ctx context.Context, id string) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
	CacheRestore(ctx context.Context) ([]*domain.Order, error)
}
