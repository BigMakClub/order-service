package usecase

import "order-service/internal/domain"

type Cache interface {
	Get(id string) (*domain.Order, bool)
	Set(id string, order *domain.Order)
}
