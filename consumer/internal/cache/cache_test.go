package cacheInMemory

import (
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"
	"testing"

	"github.com/google/uuid"
)

func TestCache_Add(t *testing.T) {
	cfg := config.New()
	cfg.Cache = config.CacheConfig{
		Capacity: 3,
	}
	c := New(cfg)

	uuids1 := []uuid.UUID{}
	for range 3 {
		uuid := uuid.New()
		uuids1 = append(uuids1, uuid)
		c.Add(models.Order{OrderUID: uuid})
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
	c.Add(models.Order{OrderUID: uuidNew})

	_, err := c.GetOrder(uuids1[0])
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
