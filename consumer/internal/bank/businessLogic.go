package bank

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/segmentio/kafka-go"
)

type Storage interface {
	// AddTransaction(t models.Transaction) error
	// GetTransaction(uuidTr uuid.UUID) (models.Transaction, error)
	// UpdateTransaction(t models.Transaction) error
	// SearchUnprocessedTransactions() []models.Transaction
}

// type CacheStorage interface {
// AddTransaction(tr models.Transaction) error
// UpdateTransaction(tr models.Transaction) error
// Clear()
// GetTransaction(uuid uuid.UUID) (models.Transaction, error)
// }

type TransactionManager struct {
	// CacheStorage CacheStorage
	// TrStorage Storage
}

func New() (*TransactionManager, error) {
	tm := &TransactionManager{
		// CacheStorage: cacheStorage,
		// TrStorage:    storage,
	}

	return tm, nil
}

func consume(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	ctx, cancel := context.WithCancel(ctx)
	var counter atomic.Int32
	ch := make(chan struct{})
	// defer close(ch)

	go func() {
		defer cancel()
		<-ch
	}()

	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{broker1Address, broker2Address, broker3Address},
				Topic:   topic,
				GroupID: "my-group",
			})
			for {
				// the `ReadMessage` method blocks until we receive the next event
				msg, err := r.ReadMessage(ctx)
				if err != nil {
					fmt.Println("could not read message " + err.Error())
					break
				}
				// after receiving the message, log its value
				fmt.Println("received: ", i, string(msg.Value))
				counter.Add(1)
				if counter.Load() == 10 {
					ch <- struct{}{}
				}
			}
			fmt.Printf("reader %v was closed\n", i)
			wg.Done()
		}(i)
	}
}
