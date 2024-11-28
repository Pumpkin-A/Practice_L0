package db

import (
	"practiceL0_go_mod/internal/models"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	Details   details   `json:"order_details"`
}

func convertToDbOrder(order models.Order) Order {
	details := details{
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          delivery(order.Delivery),
		Payment:           payment(order.Payment),
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
	for i := range order.Items {
		details.Items = append(details.Items, item(order.Items[i]))
	}
	return Order{UUID: order.OrderUID, Details: details, CreatedAt: time.Now()}
}

func convertFromDbOrder(orderTable Order) models.Order {
	order := models.Order{
		OrderUID:          orderTable.UUID,
		TrackNumber:       orderTable.Details.TrackNumber,
		Entry:             orderTable.Details.Entry,
		Delivery:          models.Delivery(orderTable.Details.Delivery),
		Payment:           models.Payment(orderTable.Details.Payment),
		Locale:            orderTable.Details.Locale,
		InternalSignature: orderTable.Details.InternalSignature,
		CustomerID:        orderTable.Details.CustomerID,
		DeliveryService:   orderTable.Details.DeliveryService,
		Shardkey:          orderTable.Details.Shardkey,
		SmID:              orderTable.Details.SmID,
		DateCreated:       orderTable.Details.DateCreated,
		OofShard:          orderTable.Details.OofShard,
	}
	for i := range orderTable.Details.Items {
		order.Items = append(order.Items, models.Item(orderTable.Details.Items[i]))
	}
	return order
}

type delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type details struct {
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          delivery  `json:"delivery"`
	Payment           payment   `json:"payment"`
	Items             []item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}
