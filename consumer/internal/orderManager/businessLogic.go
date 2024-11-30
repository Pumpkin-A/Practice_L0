package orderManager

import (
	"encoding/json"
	"log/slog"
	"practiceL0_go_mod/internal/models"

	"github.com/google/uuid"
)

type CacheStorage interface {
	AddToDBAndCache(order models.Order) error
	GetOrder(uuid uuid.UUID) (*models.Order, error)
}

type OrderManager struct {
	CacheStorage CacheStorage
}

func New(cacheStorage CacheStorage) (*OrderManager, error) {
	om := &OrderManager{
		CacheStorage: cacheStorage,
	}

	return om, nil
}

func (om *OrderManager) SaveOrder(msg []byte) error {
	order := models.Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		slog.Error("msg unmarshaling error", "func", "SaveOrder", "err", err.Error())
		return err
	}

	if !validateOrder(order) {
		slog.Error("validation error in order", "func", "SaveOrder", "order", order.OrderUID)
		return models.ErrorValidation
	}

	err = om.CacheStorage.AddToDBAndCache(order)
	if err != nil {
		slog.Error("error with adding order to DB and cache", "func", "SaveOrder", "order", order.OrderUID)
		return err
	}
	return nil
}

func (om *OrderManager) GetOrderByUUID(req models.GetOrderReq) (*models.Order, error) {
	order, err := om.CacheStorage.GetOrder(req.UUID)
	if err != nil {
		slog.Error("error with get order from cache", "func", "GetOrderByUUID", "order", order.OrderUID, "err", err.Error())
		return nil, err
	}

	slog.Info("The order with uuid: has been received", "order", order.OrderUID)
	return order, nil
}
