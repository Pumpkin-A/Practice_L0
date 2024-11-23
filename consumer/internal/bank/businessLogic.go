package bank

import (
	"context"
	"encoding/json"
	"log"
	"practiceL0_go_mod/internal/models"

	"github.com/google/uuid"
)

type Consumer interface {
	GetOrdersChan() chan []byte
	Start(ctx context.Context)
}

type CacheStorage interface {
	AddToDBAndCache(order models.Order) error
	GetOrder(uuid uuid.UUID) (*models.Order, error)
}

type TransactionManager struct {
	Consumer     Consumer
	CacheStorage CacheStorage
}

func New(consumer Consumer, cacheStorage CacheStorage) (*TransactionManager, error) {
	tm := &TransactionManager{
		Consumer:     consumer,
		CacheStorage: cacheStorage,
	}

	// tm.Consumer.Start(context.Background())
	// go tm.AddConsumedOrdersToDB()

	return tm, nil
}

func (tm *TransactionManager) AddConsumedOrdersToDB() {
	ch := tm.Consumer.GetOrdersChan()
	for {
		msg := <-ch

		order := models.Order{}
		err := json.Unmarshal(msg, &order)
		if err != nil {
			log.Println("[AddConsumedOrdersToDB] msg unmarshaling error")
		}

		// tm.Storage.Insert(order)
		err = tm.CacheStorage.AddToDBAndCache(order)
		if err != nil {
			log.Println("[AddConsumedOrdersToDB] error with adding order to DB and cache")
		}
	}
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
