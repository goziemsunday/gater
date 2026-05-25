package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrConflict          = errors.New("resource already exists")
	ErrNotFound          = errors.New("resource not found")
	queryTimeoutDuration = time.Second * 5
)

type Store struct {
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
		MarkVerified(ctx context.Context, email string) error
	}
	Verifications interface {
		Create(ctx context.Context, params CreateVerificationParams) error
		Get(ctx context.Context, hashedToken string) (*Verifications, error)
		Delete(ctx context.Context, ID string) error
	}
}

func New(pool *pgxpool.Pool) Store {
	return Store{
		Users:         &UserStore{pool},
		Verifications: &VerificationStore{pool},
	}
}
