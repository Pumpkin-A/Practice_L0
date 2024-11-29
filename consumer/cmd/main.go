package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/api"
	"practiceL0_go_mod/internal/bank"
	"practiceL0_go_mod/internal/cache"
	"practiceL0_go_mod/internal/consumer"
	"practiceL0_go_mod/internal/db"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Create context that listens for the interrupt signal from the OS.
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	pdb := db.New(cfg)
	defer func() {
		pdb.DB.Close()
		log.Println("DB was closed")
	}()

	cache := cache.New(cfg, pdb)
	tm, _ := bank.New(cache)

	consumer := consumer.New(cfg, tm)

	server, err := api.New(cfg, tm)
	if err != nil {
		log.Fatalf("Application run error: %v", err)
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		defer stop()

		if err := server.RunHTTPServer(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Application run error: %v", err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		defer log.Println("consumer was closed")
		defer stop()

		consumer.Run(mainCtx)
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		return server.Srv.Close()
	})

	// в это же время должен завершиться http-сервер, его ждём тоже...
	if err := g.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}
	log.Println("Server exiting")
}
