package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationStore struct {
	pool *pgxpool.Pool
}

type Verifications struct {
	ID         string    `json:"id"`
	Identifier string    `json:"identifier"`
	Value      string    `json:"value"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateVerificationParams struct {
	Identifier  string
	HashedToken string
	ExpiresAt   time.Time
}

func (v *VerificationStore) Create(ctx context.Context, params CreateVerificationParams) error {
	query := `
    INSERT INTO verifications (identifier, value, expires_at)
    VALUES ($1, $2, $3)
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := v.pool.Exec(
		ctx, query, params.Identifier,
		params.HashedToken, params.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create verification: %w", err)
	}

	return nil
}

func (v *VerificationStore) Get(ctx context.Context, hashedToken string) (*Verifications, error) {
	query := `
    SELECT id, identifier, value, expires_at
    FROM verifications
    WHERE value = $1
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	verification := &Verifications{}
	err := v.pool.QueryRow(ctx, query, hashedToken).Scan(
		&verification.ID, &verification.Identifier, &verification.Value, &verification.ExpiresAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return verification, nil
}

func (v *VerificationStore) Delete(ctx context.Context, ID string) error {
	query := `
    DELETE FROM verifications
    WHERE id = $1
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := v.pool.Exec(ctx, query, ID)
	if err != nil {
		return err
	}

	return nil
}
