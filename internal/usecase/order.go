package usecase

import (
	"context"
	"order-service/order-service/internal/domain"
	"order-service/order-service/internal/repository"
)

type OrderUC struct {
	repo  repository.OrderRepository
	cache Cache
}

func NewOrederUC(r repository.OrderRepository, c Cache) *OrderUC {
	return &OrderUC{r, c}
}

func (uc *OrderUC) Get(ctx context.Context, id string) (*domain.Order, error) {
	if v, ok := uc.cache.Get(id); ok {
		return v, nil
	}

	order, err := uc.repo.Find(ctx, id)

	if err == nil {
		uc.cache.Set(id, order)
	}

	return order, nil
}

func (uc *OrderUC) Set(ctx context.Context, order *domain.Order) error {
	if err := uc.repo.Save(ctx, order); err != nil {
		return err
	}
	uc.cache.Set(order.OrderId.String(), order)
	return nil
}
