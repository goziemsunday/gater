package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/goziemsunday/gater/internal/config"
	"github.com/goziemsunday/gater/internal/db"
	"github.com/goziemsunday/gater/internal/mailer"
	"github.com/goziemsunday/gater/internal/store"
	"github.com/goziemsunday/gater/internal/validator"
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
	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		logger.Error("failed to create db pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("database connection pool established")

	dbStore := store.New(pool)

	emailer := mailer.NewResendClient(cfg)

	validator := validator.New()

	// init app
	app := &application{
		config:    cfg,
		store:     dbStore,
		mailer:    emailer,
		validator: validator,
		logger:    logger,
	}

	if err := app.run(app.mount()); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
