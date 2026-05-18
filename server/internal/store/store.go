package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var queryTimeoutDuration = time.Second * 5

type Store struct {
	Users interface {
		Create(ctx context.Context, params CreateUserParams) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
	}
}

func New(pool *pgxpool.Pool) Store {
	return Store{
		Users: &UserStore{pool},
	}
}
