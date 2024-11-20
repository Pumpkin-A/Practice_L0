package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
	"practiceL0_go_mod/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
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
func (d *details) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &d)
}

func (pdb *PostgresDB) Insert(order models.Order) {
	orderTable := convertToDbDetails(order)
	_, err := pdb.DB.Exec("INSERT INTO orders (uuid, details) VALUES($1, $2)", orderTable.UUID, orderTable.Details)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("данные успешно записаны")
}

func (pdb *PostgresDB) SelectByUUID(uuid uuid.UUID) (*models.Order, error) {
	orderInDB := &OrderTable{}
	err := pdb.DB.QueryRow("SELECT uuid, details FROM orders WHERE uuid = $1", uuid).Scan(&orderInDB.UUID, &orderInDB.Details)
	if err != nil {
		log.Println("[SelectByUUID] error with get order from db")
		return nil, err
	}
	order := convertToOrder(*orderInDB)
	return &order, nil
}
