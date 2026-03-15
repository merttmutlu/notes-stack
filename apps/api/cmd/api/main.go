package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/merttmutlu/notes-stack/apps/api/internal/api"
	"github.com/merttmutlu/notes-stack/apps/api/internal/config"
	"github.com/merttmutlu/notes-stack/apps/api/internal/store"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	store, err := store.New(context.Background(), store.Dependencies{
		DatabaseURL: cfg.DatabaseURL,
	})
	if err != nil {
		logger.Error("failed to initialize store", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	router := api.NewRouter(api.Dependencies{
		Logger: logger,
		Store:  store,
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	logger.Info("starting api server", "port", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}

	_ = server.Shutdown(context.Background())
}
