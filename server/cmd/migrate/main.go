package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chiagxziem/snipper/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/pflag"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// get path to migrations from flag
	pathFlag := pflag.StringP("path", "p", "", "path to migrations folder")
	pflag.Parse()
	if *pathFlag == "" {
		log.Fatal("path flag is required")
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+*pathFlag, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("no migrations to apply")
	} else {
		fmt.Println("migrations successful!")
	}

}
