package bank

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
