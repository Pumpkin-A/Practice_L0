package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/api"
	"practiceL0_go_mod/internal/cache"
	"practiceL0_go_mod/internal/consumer"
	"practiceL0_go_mod/internal/db"
	"practiceL0_go_mod/internal/orderManager"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create context that listens for the interrupt signal from the OS.
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	pdb := db.New(cfg)
	defer func() {
		pdb.DB.Close()
		slog.Info("DB was closed")
	}()

	cache := cache.New(cfg, pdb)
	om, _ := orderManager.New(cache)

	consumer := consumer.New(cfg, om)

	server, err := api.New(cfg, om)
	if err != nil {
		slog.Error("Application run error", "err", err.Error())
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		defer stop()

		if err := server.RunHTTPServer(); err != nil && err != http.ErrServerClosed {
			slog.Error("Application run error", "err", err.Error())
			return err
		}
		return nil
	})
	g.Go(func() error {
		defer slog.Info("consumer was closed")
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
		slog.Info("server exit reason", "err", err.Error())
	}
	slog.Info("Server exiting")
}
