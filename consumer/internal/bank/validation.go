package bank

import (
	"log"
	"practiceL0_go_mod/internal/models"
	"time"

	"github.com/google/uuid"
)

func validateItem(item models.Item) bool {
	if item.ChrtID <= 0 || item.NmID <= 0 {
		return false
	}
	if item.Price <= 0 || item.Sale < 0 || item.TotalPrice <= 0 ||
		item.TotalPrice < item.Price || item.TotalPrice < item.Sale {
		log.Println("invalid price or sale")
		return false
	}
	if item.Rid == "" || item.Name == "" || item.Size == "" || item.Brand == "" || item.Status == 0 {
		return false
	}
	return true
}

func validatePayment(payment models.Payment) bool {
	var nilUID uuid.UUID
	if payment.RequestID == nilUID.String() || payment.Transaction == nilUID.String() {
		return false
	}
	if payment.Currency == "" {
		return false
	}
	if payment.Provider != "wbpay" {
		return false
	}
	if payment.PaymentDt > time.Now().Unix() {
		return false
	}
	if payment.Amount <= 0 || payment.DeliveryCost < 0 || payment.GoodsTotal <= 0 {
		return false
	}
	if payment.Bank == "" {
		return false
	}
	return true
}

func validateDelivery(delivery models.Delivery) bool {
	if delivery.Address == "" || delivery.City == "" || delivery.Email == "" || delivery.Name == "" ||
		delivery.Phone == "" || delivery.Region == "" || delivery.Zip == "" {
		return false
	}
	return true
}

func validateOrder(order models.Order) bool {
	var nilUID uuid.UUID
	if order.OrderUID == nilUID {
		return false
	}
	if order.TrackNumber == "" || order.Entry == "" {
		return false
	}
	if !validateDelivery(order.Delivery) {
		return false
	}
	if !validatePayment(order.Payment) {
		return false
	}
	if order.Locale == "" || order.CustomerID == "" || order.DeliveryService == "" {
		return false
	}
	if order.SmID <= 0 {
		return false
	}
	if order.DateCreated.After(time.Now()) {
		return false
	}
	if len(order.Items) <= 0 {
		return false
	}
	for _, item := range order.Items {
		if !validateItem(item) {
			return false
		}
	}

	return true
}
