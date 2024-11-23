package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func New(cfg config.Config) *PostgresDB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", cfg.DB.DbUser, cfg.DB.DbPassword, cfg.DB.DbName, cfg.DB.SSLmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connection success")

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

func (pdb *PostgresDB) Insert(order models.Order) error {
	orderTable := convertToDbOrder(order)
	_, err := pdb.DB.Exec("INSERT INTO orders (uuid, details, created_at) VALUES($1, $2, $3)", orderTable.UUID, orderTable.Details, orderTable.CreatedAt)
	if err != nil {
		log.Printf("[SelectByUUID] error with adding order to DB with uuid: %v %s\n", order.OrderUID, err.Error())
		return err
	}
	log.Println("order was successfully added to DB")
	return nil
}

func (pdb *PostgresDB) GetOrderByUUID(uuid uuid.UUID) (*models.Order, error) {
	orderInDB := &Order{}
	err := pdb.DB.QueryRow("SELECT uuid, details FROM orders WHERE uuid = $1", uuid).Scan(&orderInDB.UUID, &orderInDB.Details)
	if err != nil {
		log.Println("[SelectByUUID] error with get order from db", err.Error())
		return nil, err
	}
	order := convertFromDbOrder(*orderInDB)
	return &order, nil
}

func (pdb *PostgresDB) CacheRecovery(limit int) ([]models.Order, error) {
	ordersInDB := []*Order{}
	rows, err := pdb.DB.Query("SELECT * FROM (SELECT * FROM orders o ORDER BY o.created_at DESC LIMIT $1) AS tbl ORDER BY tbl.created_at ASC;", limit)
	if err != nil {
		log.Println("[CacheRecovery] error with get orders from db", err.Error())
		return nil, err
	}
	for rows.Next() {
		orderInDB := Order{}
		err := rows.Scan(&orderInDB.UUID, &orderInDB.Details, &orderInDB.CreatedAt)
		if err != nil {
			log.Println("[CacheRecovery] error with scanning ordersInDB to orders", err.Error())
			return nil, err
		}
		ordersInDB = append(ordersInDB, &orderInDB)
	}

	orders := []models.Order{}
	for _, orderInDB := range ordersInDB {
		order := convertFromDbOrder(*orderInDB)
		orders = append(orders, order)
	}

	return orders, nil
}
