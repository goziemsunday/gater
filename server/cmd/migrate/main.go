package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/chiagxziem/snipper/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/pflag"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// get path to migrations from flag
	pathFlag := pflag.StringP("path", "p", "", "path to migrations folder")
	pflag.Parse()
	if *pathFlag == "" {
		fmt.Fprintln(os.Stderr, "path flag is required")
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+*pathFlag, "postgres", driver)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("no migrations to apply")
	} else {
		fmt.Println("migrations successful!")
	}

}
