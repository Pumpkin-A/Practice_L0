package cache

import (
	"log"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"
	"sync"

	"github.com/google/uuid"
)

type Storage interface {
	Insert(order models.Order) error
	GetOrderByUUID(uuid uuid.UUID) (*models.Order, error)
	CacheRecovery(limit int) ([]models.Order, error)
}

type Cache struct {
	storage       Storage
	mutex         sync.RWMutex
	items         map[uuid.UUID]models.Order
	capacity      int
	mainBuf       []uuid.UUID
	additionalBuf []uuid.UUID
	index         int
}

func New(cfg config.Config, storage Storage) *Cache {
	items := make(map[uuid.UUID]models.Order)
	mainBuf := make([]uuid.UUID, cfg.Cache.Capacity)
	additionalBuf := make([]uuid.UUID, cfg.Cache.Capacity)
	cache := Cache{
		storage:       storage,
		items:         items,
		capacity:      cfg.Cache.Capacity,
		mainBuf:       mainBuf,
		additionalBuf: additionalBuf,
	}

	go cache.recovery()

	return &cache
}

func (c *Cache) add(order models.Order) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.index == c.capacity {
		c.mainBuf, c.additionalBuf = c.additionalBuf, c.mainBuf
		c.index = 0
	}

	var nilUUID uuid.UUID
	delete(c.items, c.additionalBuf[c.index])
	c.additionalBuf[c.index] = nilUUID

	c.items[order.OrderUID] = order
	c.mainBuf[c.index] = order.OrderUID
	c.index++

	return nil
}

func (c *Cache) GetOrder(uuid uuid.UUID) (*models.Order, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	order, ok := c.items[uuid]
	if !ok {
		log.Printf("[cache GetOrder] order with uuid %v was not found in cache\n", uuid)

		order, err := c.storage.GetOrderByUUID(uuid)
		if err != nil {
			log.Printf("[GetOrderByUUID] error with get order from db: %s", err.Error())
			return nil, models.ErrorOrderNotExist
		}
		return order, nil
	}
	return &order, nil
}

func (c *Cache) AddToDBAndCache(order models.Order) error {
	err := c.storage.Insert(order)
	if err != nil {
		return err
	}

	err = c.add(order)
	if err != nil {
		return err
	}

	log.Printf("Order with uuid: %v was successfully added to DB and cache\n", order.OrderUID)
	return nil
}

func (c *Cache) recovery() error {
	orders, err := c.storage.CacheRecovery(c.capacity)
	if err != nil {
		return err
	}

	for _, order := range orders {
		err = c.add(order)
		if err != nil {
			return err
		}
	}
	return nil
}
