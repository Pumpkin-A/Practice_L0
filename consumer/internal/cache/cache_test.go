package cache

import (
	"errors"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"
	"testing"

	"github.com/google/uuid"
)

type mapStorage struct {
	Orders map[uuid.UUID]models.Order
}

func NewMap() mapStorage {
	return mapStorage{
		Orders: make(map[uuid.UUID]models.Order),
	}
}

func (ms *mapStorage) Insert(order models.Order) error {
	ms.Orders[order.OrderUID] = order
	return nil
}

func (ms *mapStorage) GetOrderByUUID(uuid uuid.UUID) (*models.Order, error) {
	order, ok := ms.Orders[uuid]
	if !ok {
		return nil, errors.New("order not exist")
	}
	return &order, nil
}

func (ms *mapStorage) CacheRecovery(limit int) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	var i int
	for _, value := range ms.Orders {
		// для теста реализовано восстановление без учета даты создания заказа
		if i >= limit {
			break
		}
		orders = append(orders, value)
		i++
	}
	return orders, nil
}

func TestCache_Recovery(t *testing.T) {
	cfg := config.New()
	cfg.Cache = config.CacheConfig{
		Capacity: 3,
	}
	storage := NewMap()
	c := New(cfg, &storage)

	uuids := []uuid.UUID{}
	for range 3 {
		uuid := uuid.New()
		uuids = append(uuids, uuid)
		c.storage.Insert(models.Order{OrderUID: uuid})
	}

	c.recovery()

	{
		_, err := c.GetOrder(uuids[0])
		if err != nil {
			t.Error(err)
		}
	}
	{
		_, err := c.GetOrder(uuids[1])
		if err != nil {
			t.Error(err)
		}
	}
	{
		_, err := c.GetOrder(uuids[2])
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCache_Add(t *testing.T) {
	cfg := config.New()
	cfg.Cache = config.CacheConfig{
		Capacity: 3,
	}
	storage := NewMap()
	c := New(cfg, &storage)

	uuids1 := []uuid.UUID{}
	for range 3 {
		uuid := uuid.New()
		uuids1 = append(uuids1, uuid)
		c.AddToDBAndCache(models.Order{OrderUID: uuid})
	}

	for _, uuid := range uuids1 {
		got, err := c.GetOrder(uuid)
		if err != nil {
			t.Error(err)
		}
		if got.OrderUID != uuid {
			t.Errorf("uuids not equal")
		}
	}

	uuidNew := uuid.New()
	c.AddToDBAndCache(models.Order{OrderUID: uuidNew})

	_, err := c.GetOrder(uuids1[0])
	if err != nil {
		t.Error(err)
	}

	delete(storage.Orders, uuids1[0])

	_, err = c.GetOrder(uuids1[0])
	if nil == err {
		t.Error("error expected")
	}
	{
		_, err := c.GetOrder(uuids1[1])
		if err != nil {
			t.Error(err)
		}
	}
	{
		_, err := c.GetOrder(uuids1[2])
		if err != nil {
			t.Error(err)
		}
	}
	{
		_, err := c.GetOrder(uuidNew)
		if err != nil {
			t.Error(err)
		}
	}
}
