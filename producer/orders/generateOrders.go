package orders

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func GenerateOrder() Order {
	order := Order{
		OrderUID:    uuid.New(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  gofakeit.TimeZoneRegion(),
			Email:   gofakeit.Email(),
		},
		Payment: Payment{
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
		Items:             []Item{},
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
		item := Item{
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
