package cache

import (
	"order-service/order-service/internal/domain"
	"sync"
)

type MemCache struct {
	mu   sync.RWMutex
	data map[string]*domain.Order
}

func NewMemCache(capacity int) *MemCache {
	return &MemCache{
		data: make(map[string]*domain.Order, capacity),
	}
}

func (c *MemCache) Get(id string) (*domain.Order, bool) {
	c.mu.RLock()
	order, ok := c.data[id]
	c.mu.RUnlock()
	return order, ok
}

func (c *MemCache) Set(id string, order *domain.Order) {
	c.mu.Lock()
	c.data[id] = order
	c.mu.Unlock()
}
