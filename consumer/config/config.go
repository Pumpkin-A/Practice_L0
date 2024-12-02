package config

import "os"

type DBConfig struct {
	DbHost     string
	DbPort     string
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
type CacheConfig struct {
	Capacity int
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Kafka  KafkaConfig
	Cache  CacheConfig
}

func New() Config {
	db := DBConfig{
		SSLmode: "disable",
	}
	DbName, exists := os.LookupEnv("POSTGRES_DB")
	if exists {
		db.DbName = DbName
	}
	DbHost, exists := os.LookupEnv("POSTGRES_HOST")
	if exists {
		db.DbHost = DbHost
	}
	DbPort, exists := os.LookupEnv("POSTGRES_PORT")
	if exists {
		db.DbPort = DbPort
	}
	DbUser, exists := os.LookupEnv("POSTGRES_USER")
	if exists {
		db.DbUser = DbUser
	}
	dbPassword, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if exists {
		db.DbPassword = dbPassword
	}

	kafka := KafkaConfig{
		Topic:             "orders",
		Broker2Address:    "localhost:9093",
		Broker3Address:    "localhost:9094",
		NumberOfConsumers: 5,
	}

	broker, exists := os.LookupEnv("KAFKA_BROKER")
	if exists {
		kafka.Broker1Address = broker
	}

	server := ServerConfig{
		Port: 9090,
	}

	cache := CacheConfig{
		Capacity: 100,
	}

	return Config{DB: db,
		Server: server,
		Kafka:  kafka,
		Cache:  cache}
}
