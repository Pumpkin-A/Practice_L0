package orderManager

import (
	"encoding/json"
	"math/rand"
	"practiceL0_go_mod/internal/models"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type CacheStorageOk struct {
}

func CacheNew() *CacheStorageOk {
	return &CacheStorageOk{}
}

func (cs *CacheStorageOk) AddToDBAndCache(order models.Order) error {
	return nil
}

func (cs *CacheStorageOk) GetOrder(uuid uuid.UUID) (*models.Order, error) {
	return nil, nil
}

func TestOrderManager_Validation(t *testing.T) {
	cache := CacheNew()
	om := New(cache)

	fakeOrder := &fakeOrder{
		uuid: uuid.New(),
	}
	fakeOrderBytes, _ := json.Marshal(fakeOrder)
	err := om.SaveOrder(fakeOrderBytes)
	if err == nil {
		t.Error("error expected")
	}

	order1 := generateOrder()
	order1Bytes, _ := json.Marshal(order1)
	err = om.SaveOrder(order1Bytes)
	if err != nil {
		t.Error("error with unmarshal")
	}

	order2 := generateOrder()
	order2.Delivery.Name = ""
	order2.Items[0].Sale = 300
	order2.Items[0].TotalPrice = 200

	order2Bytes, _ := json.Marshal(order2)
	err = om.SaveOrder(order2Bytes)
	if err != models.ErrorValidation {
		t.Error("validation error expected")
	}

}

type fakeOrder struct {
	uuid uuid.UUID
}

func generateOrder() models.Order {
	order := models.Order{
		OrderUID:    uuid.New(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  gofakeit.TimeZoneRegion(),
			Email:   gofakeit.Email(),
		},
		Payment: models.Payment{
			Transaction:  gofakeit.UUID(),
			RequestID:    gofakeit.UUID(),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     "wbpay",
			Amount:       gofakeit.IntN(500000) + 1,
			PaymentDt:    gofakeit.DateRange(time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC), time.Now()).Unix(),
			Bank:         gofakeit.RandomString([]string{"alpha", "sberBank", "t-bank", "vtb"}),
			DeliveryCost: gofakeit.IntN(10000),
			GoodsTotal:   gofakeit.IntN(1000) + 1,
			CustomFee:    0,
		},
		Items:             []models.Item{},
		Locale:            gofakeit.RandomString([]string{"en", "ru"}),
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              gofakeit.IntN(300) + 1,
		DateCreated:       gofakeit.DateRange(time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC), time.Now()),
		OofShard:          "1",
	}

	itemsCount := rand.Intn(10) + 1
	for range itemsCount {
		item := models.Item{
			ChrtID:      gofakeit.IntN(1000000),
			TrackNumber: "WBILMTESTTRACK",
			Price:       gofakeit.IntN(50000),
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        gofakeit.IntN(500),
			Size:        "0",
			NmID:        gofakeit.IntN(1000000),
			Brand:       gofakeit.RandomString([]string{"Vivienne Sabo", "Prada", "Gucci"}),
			Status:      gofakeit.RandomInt([]int{202, 200, 400}),
		}
		item.TotalPrice = item.Price + item.Sale

		order.Items = append(order.Items, item)
	}

	return order
}
