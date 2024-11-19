package main

import (
	"log"
	"practiceL0_go_mod/internal/api"
	"practiceL0_go_mod/internal/bank"
	"practiceL0_go_mod/internal/consumer"
	"practiceL0_go_mod/internal/db"
)

func main() {
	// diskStorage, err := diskTrStorage.New()
	// defer diskStorage.Close()
	// if err != nil {
	// 	log.Fatalf("Application run error: %v", err)
	// }
	pdb := db.New("user=postgres password=mypass dbname=ordersdb sslmode=disable")
	consumer := consumer.New()
	tm, _ := bank.New(consumer, pdb)
	server, err := api.New(tm)
	if err != nil {
		log.Fatalf("Application run error: %v", err)
	}
	err = server.Run()
	if err != nil {
		log.Fatalf("Application run error: %v", err)
	}
	log.Println("Application was ran successfully")
}
