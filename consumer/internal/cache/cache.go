package cacheInMemory

import (
	"errors"
	"log"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mutex         sync.RWMutex
	items         map[uuid.UUID]models.Order
	capacity      int
	mainBuf       []uuid.UUID
	additionalBuf []uuid.UUID
	index         int
}

func New(cfg config.Config) *Cache {
	items := make(map[uuid.UUID]models.Order)
	mainBuf := make([]uuid.UUID, cfg.Cache.Capacity)
	additionalBuf := make([]uuid.UUID, cfg.Cache.Capacity)
	cache := Cache{
		items:         items,
		capacity:      cfg.Cache.Capacity,
		mainBuf:       mainBuf,
		additionalBuf: additionalBuf,
	}

	return &cache
}

func (c *Cache) Add(order models.Order) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var nilUUID uuid.UUID
	if c.additionalBuf[c.index] == nilUUID {
		c.items[order.OrderUID] = order
		c.mainBuf[c.index] = order.OrderUID
	} else {
		delete(c.items, c.additionalBuf[c.index])
		c.additionalBuf[c.index] = nilUUID
		c.items[order.OrderUID] = order
		c.mainBuf[c.index] = order.OrderUID
	}
	c.index++

	if c.index == c.capacity {
		buf := c.additionalBuf
		c.additionalBuf = c.mainBuf
		c.mainBuf = buf
		c.index = 0
	}
	return nil
}

func (c *Cache) GetOrder(uuid uuid.UUID) (*models.Order, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	order, ok := c.items[uuid]
	if !ok {
		log.Printf("[cache GetOrder] order with uuid %v was not found in cache\n", uuid)
		return nil, errors.New("order was not found in cache")
	}
	return &order, nil
}
