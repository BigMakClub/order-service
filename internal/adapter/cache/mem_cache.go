package cache

import (
	"log"
	"order-service/internal/domain"
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
	defer c.mu.RUnlock()
	
	log.Printf("Cache GET request for ID: %s", id)
	order, ok := c.data[id]
	if !ok {
		log.Printf("Cache MISS for ID: %s", id)
		return nil, false
	}
	log.Printf("Cache HIT for ID: %s", id)
	return order, true
}

func (c *MemCache) Set(id string, order *domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	log.Printf("Cache SET for ID: %s", id)
	c.data[id] = order
	log.Printf("Cache now contains %d items", len(c.data))
}
