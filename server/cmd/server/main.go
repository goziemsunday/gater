package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/chiagxziem/snipper/internal/config"
	"github.com/chiagxziem/snipper/internal/db"
	"github.com/chiagxziem/snipper/internal/store"
	"github.com/chiagxziem/snipper/internal/validator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// load app config
	cfg, err := config.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// database
	pool, err := db.NewPool(ctx)
	if err != nil {
		logger.Error("failed to create db pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("database connection pool established")

	// init dbStore
	dbStore := store.New(pool)

	validator := validator.New()

	// init app
	app := &application{
		config:    cfg,
		store:     dbStore,
		validator: validator,
		logger:    logger,
	}

	if err := app.run(app.mount()); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
