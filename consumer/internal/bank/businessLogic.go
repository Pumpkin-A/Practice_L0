package bank

import (
	"context"
)

type Storage interface {
	// AddTransaction(t models.Transaction) error
	// GetTransaction(uuidTr uuid.UUID) (models.Transaction, error)
	// UpdateTransaction(t models.Transaction) error
	// SearchUnprocessedTransactions() []models.Transaction
}

type Consumer interface {
	Consume(ctx context.Context)
}

// type CacheStorage interface {
// AddTransaction(tr models.Transaction) error
// UpdateTransaction(tr models.Transaction) error
// Clear()
// GetTransaction(uuid uuid.UUID) (models.Transaction, error)
// }

type TransactionManager struct {
	Consumer Consumer
	// CacheStorage CacheStorage
	// TrStorage Storage
}

func New(consumer Consumer) (*TransactionManager, error) {
	tm := &TransactionManager{
		Consumer: consumer,
		// CacheStorage: cacheStorage,
		// TrStorage:    storage,
	}

	tm.Consumer.Consume(context.Background())

	return tm, nil
}
