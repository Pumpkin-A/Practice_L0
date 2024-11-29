package bank

import (
	"encoding/json"
	"log"
	"practiceL0_go_mod/internal/models"

	"github.com/google/uuid"
)

type CacheStorage interface {
	AddToDBAndCache(order models.Order) error
	GetOrder(uuid uuid.UUID) (*models.Order, error)
}

type TransactionManager struct {
	CacheStorage CacheStorage
}

func New(cacheStorage CacheStorage) (*TransactionManager, error) {
	tm := &TransactionManager{
		CacheStorage: cacheStorage,
	}

	return tm, nil
}

func (tm *TransactionManager) AddConsumedOrdersToDBAndCache(msg []byte) error {
	order := models.Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Println("[AddConsumedOrdersToDBAndCache] msg unmarshaling error")
		return err
	}

	if !validateOrder(order) {
		log.Printf("[AddConsumedOrdersToDBAndCache] validation error in order: %v\n", order)
		return models.ErrorValidation
	}

	err = tm.CacheStorage.AddToDBAndCache(order)
	if err != nil {
		log.Println("[AddConsumedOrdersToDBAndCache] error with adding order to DB and cache")
		return err
	}
	return nil
}

func (tm *TransactionManager) GetOrderByUUID(req models.GetOrderReq) (*models.Order, error) {
	order, err := tm.CacheStorage.GetOrder(req.UUID)
	if err != nil {
		log.Printf("[GetOrderByUUID] error with get order from cache: %s", err.Error())
		return nil, err
	}

	log.Printf("The order with uuid: %v has been received\n", order.OrderUID)
	return order, nil
}
