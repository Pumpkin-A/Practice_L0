package main

import (
	"log"
	"practiceL0_go_mod/internal/api"
	"practiceL0_go_mod/internal/bank"
)

func main() {
	// diskStorage, err := diskTrStorage.New()
	// defer diskStorage.Close()
	// if err != nil {
	// 	log.Fatalf("Application run error: %v", err)
	// }
	tm, _ := bank.New()
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
