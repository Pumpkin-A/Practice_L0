package main

import (
	"log"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/api"
	"practiceL0_go_mod/internal/bank"
	"practiceL0_go_mod/internal/cache"
	"practiceL0_go_mod/internal/consumer"
	"practiceL0_go_mod/internal/db"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.New()
	pdb := db.New(cfg)
	cache := cache.New(cfg, pdb)
	consumer := consumer.New(cfg)
	tm, _ := bank.New(consumer, cache)
	server, err := api.New(tm)
	if err != nil {
		log.Fatalf("Application run error: %v", err)
	}
	err = server.Run(cfg)
	if err != nil {
		log.Fatalf("Application run error: %v", err)
	}
	log.Println("Application was ran successfully")
}
