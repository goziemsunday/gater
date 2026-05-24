package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	pool *pgxpool.Pool
}

type User struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PasswordHash  *string   `json:"-"`
	EmailVerified bool      `json:"email_verified"`
	Image         *string   `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
    INSERT INTO users (name, email, password_hash, image)
    VALUES ($1, $2, $3, $4)
    RETURNING id, name, email, email_verified, image, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	err := s.pool.QueryRow(
		ctx, query, user.Name, user.Email,
		user.PasswordHash, user.Image,
	).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerified,
		&user.Image, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key"):
			return ErrEmailAlreadyExists
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
    SELECT id, name, email, email_verified, image, created_at, updated_at 
    FROM users 
    WHERE email = $1
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var user User
	err := s.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerified,
		&user.Image, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
