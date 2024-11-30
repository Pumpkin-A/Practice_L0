package orderManager

import (
	"log/slog"
	"practiceL0_go_mod/internal/models"
	"time"

	"github.com/google/uuid"
)

func validateItem(item models.Item) bool {
	if item.ChrtID <= 0 || item.NmID <= 0 {
		slog.Error("invalid ChrtID or NmID", "func", "validateItem")
		return false
	}
	if item.Price <= 0 || item.Sale < 0 || item.TotalPrice <= 0 ||
		item.TotalPrice < item.Price || item.TotalPrice < item.Sale {
		slog.Error("invalid price or sale", "func", "validateItem")
		return false
	}
	if item.Rid == "" || item.Name == "" || item.Size == "" || item.Brand == "" || item.Status == 0 {
		slog.Error("required fields are not filled in", "func", "validateItem")
		return false
	}
	return true
}

func validatePayment(payment models.Payment) bool {
	var nilUID uuid.UUID
	if payment.RequestID == nilUID.String() || payment.Transaction == nilUID.String() {
		slog.Error("required fields are not filled in", "func", "validatePayment")
		return false
	}
	if payment.Currency == "" {
		slog.Error("required fields are not filled in", "func", "validatePayment")
		return false
	}
	if payment.Provider != "wbpay" {
		slog.Error("invalid Provider", "func", "validatePayment")
		return false
	}
	if payment.PaymentDt > time.Now().Unix() {
		slog.Error("invalid PaymentDt", "func", "validatePayment")
		return false
	}
	if payment.Amount <= 0 || payment.DeliveryCost < 0 || payment.GoodsTotal <= 0 {
		slog.Error("invalid information about the cost of the product", "func", "validatePayment")
		return false
	}
	if payment.Bank == "" {
		slog.Error("invalid bank data", "func", "validatePayment")
		return false
	}
	return true
}

func validateDelivery(delivery models.Delivery) bool {
	if delivery.Address == "" || delivery.City == "" || delivery.Email == "" || delivery.Name == "" ||
		delivery.Phone == "" || delivery.Region == "" || delivery.Zip == "" {
		slog.Error("required fields are not filled in", "func", "validateDelivery")
		return false
	}
	return true
}

func validateOrder(order models.Order) bool {
	var nilUID uuid.UUID
	if order.OrderUID == nilUID {
		slog.Error("required fields are not filled in", "func", "validateOrder")
		return false
	}
	if order.TrackNumber == "" || order.Entry == "" {
		slog.Error("required fields are not filled in", "func", "validateOrder")
		return false
	}
	if !validateDelivery(order.Delivery) {
		return false
	}
	if !validatePayment(order.Payment) {
		return false
	}
	if order.Locale == "" || order.CustomerID == "" || order.DeliveryService == "" {
		slog.Error("required fields are not filled in", "func", "validateOrder")
		return false
	}
	if order.SmID <= 0 {
		slog.Error("invalid SmID", "func", "validateOrder")
		return false
	}
	if order.DateCreated.After(time.Now()) {
		slog.Error("invalid DateCreated", "func", "validateOrder")
		return false
	}
	if len(order.Items) <= 0 {
		slog.Error("invalid тumber of products", "func", "validateOrder")
		return false
	}
	for _, item := range order.Items {
		if !validateItem(item) {
			return false
		}
	}

	return true
}
