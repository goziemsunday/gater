package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chiagxziem/snipper/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := &api{
		config: cfg,
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
		log.Printf("server failed to start, err: %v", err)
		os.Exit(1)
	}
}
