package config

import "os"

type DBConfig struct {
	DbUser     string
	DbPassword string
	DbName     string
	SSLmode    string
}

type ServerConfig struct {
	Port int
}

type KafkaConfig struct {
	Topic             string
	Broker1Address    string
	Broker2Address    string
	Broker3Address    string
	NumberOfConsumers int
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Kafka  KafkaConfig
}

func New() Config {
	db := DBConfig{
		DbUser:  "postgres",
		DbName:  "ordersdb",
		SSLmode: "disable",
	}
	dbPassword, exists := os.LookupEnv("DbPassword")
	if exists {
		db.DbPassword = dbPassword
	}

	kafka := KafkaConfig{
		Topic:             "orders",
		Broker1Address:    "localhost:9092",
		Broker2Address:    "localhost:9093",
		Broker3Address:    "localhost:9094",
		NumberOfConsumers: 5,
	}

	server := ServerConfig{
		Port: 9090,
	}

	return Config{DB: db,
		Server: server,
		Kafka:  kafka}
}
