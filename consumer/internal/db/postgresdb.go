package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
	"practiceL0_go_mod/internal/models"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB          *sql.DB
	OrdersTable OrderTable
}

type OrderTable struct {
	UUID    uuid.UUID `json:"uuid"`
	Details details   `json:"order_details"`
}

func New(connStr string) *PostgresDB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success")

	return &PostgresDB{
		DB: db,
	}
}

// Valuer interface. This method simply returns the JSON-encoded representation of the struct.
func (d details) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scanner interface. This method simply decodes a JSON-encoded value into the struct fields.
func (d details) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &d)
}

func (pdb *PostgresDB) Insert(order models.Order) {
	details := convertToDbDetails(order)
	_, err := pdb.DB.Exec("INSERT INTO orders (uuid, details) VALUES($1, $2)", order.OrderUID, details)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("данные успешно записаны")
}

// func main() {
// db, err := sql.Open("postgres", "postgres://user:pass@localhost/db")
// if err != nil {
// 	log.Fatal(err)
// }

// // Initialize a new Attrs struct and add some values.
// attrs := new(Attrs)
// attrs.Name = "Pesto"
// attrs.Ingredients = []string{"Basil", "Garlic", "Parmesan", "Pine nuts", "Olive oil"}
// attrs.Organic = false
// attrs.Dimensions.Weight = 100.00

// // The database driver will call the Value() method and and marshall the
// // attrs struct to JSON before the INSERT.
// _, err = db.Exec("INSERT INTO items (attrs) VALUES($1)", attrs)
// if err != nil {
// 	log.Fatal(err)
// }

// // Similarly, we can also fetch data from the database, and the driver
// // will call the Scan() method to unmarshal the data to an Attr struct.
// item := new(Item)
// err = db.QueryRow("SELECT id, attrs FROM items ORDER BY id DESC LIMIT 1").Scan(&item.ID, &item.Attrs)
// if err != nil {
// 	log.Fatal(err)
// }

// // You can then use the struct fields as normal...
// weightKg := item.Attrs.Dimensions.Weight / 1000
// log.Printf("Item: %d, Name: %s, Weight: %.2fkg", item.ID, item.Attrs.Name, weightKg)
// }

func convertToDbDetails(order models.Order) details {
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
	return details
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
	PaymentDt    int    `json:"payment_dt"`
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
