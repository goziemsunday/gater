package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

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

	srv := &http.Server{
		Addr:         ":" + app.config.Port,
		Handler:      app.mount(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr :%s", app.config.Port)

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
