package bank

import (
	"context"
	"encoding/json"
	"log"
	"practiceL0_go_mod/internal/models"
)

type Storage interface {
	Insert(order models.Order)
}

type Consumer interface {
	GetOrdersChan() chan []byte
	Start(ctx context.Context)
}

// type CacheStorage interface {
// AddTransaction(tr models.Transaction) error
// UpdateTransaction(tr models.Transaction) error
// Clear()
// GetTransaction(uuid uuid.UUID) (models.Transaction, error)
// }

type TransactionManager struct {
	Consumer Consumer
	Storage  Storage
	// CacheStorage CacheStorage
}

func New(consumer Consumer, storage Storage) (*TransactionManager, error) {
	tm := &TransactionManager{
		Consumer: consumer,
		Storage:  storage,
		// CacheStorage: cacheStorage,
	}

	tm.Consumer.Start(context.Background())
	go tm.AddConsumedOrdersToDB()

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

		tm.Storage.Insert(order)
	}
}
