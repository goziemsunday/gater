package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var queryTimeoutDuration = time.Second * 5

type Store struct {
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
	}
	Verifications interface {
		Create(ctx context.Context, params CreateVerificationParams) error
	}
}

func New(pool *pgxpool.Pool) Store {
	return Store{
		Users:         &UserStore{pool},
		Verifications: &VerificationStore{pool},
	}
}
