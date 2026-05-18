package main

import (
	"log/slog"
	"os"

	"github.com/chiagxziem/snipper/internal/config"
)

func main() {
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// load app config
	cfg, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// init app
	app := &application{
		config: cfg,
		logger: logger,
	}

	if err := app.run(app.mount()); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
