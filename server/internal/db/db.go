package db

import (
	"context"
	"fmt"

	"github.com/chiagxziem/snipper/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// get dbCfg struct
	dbCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// create pool from config
	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// ping database
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
