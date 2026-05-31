package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OAuthStore struct {
	pool *pgxpool.Pool
}

type OAuthAccount struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	Provider              string    `json:"provider"`
	ProviderAccountID     string    `json:"provider_account_id"`
	AccessToken           string    `json:"-"`
	RefreshToken          string    `json:"-"`
	IDToken               string    `json:"-"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	Scope                 string    `json:"scope"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (s *OAuthStore) GetByProviderAndAccountID(
	ctx context.Context,
	provider, accountID string,
) (*OAuthAccount, error) {
	query := `
    SELECT id, user_id, provider, provider_account_id, access_token, refresh_token, id_token, access_token_expires_at, refresh_token_expires_at, scope, created_at, updated_at
    FROM oauth_accounts
    WHERE provider = $1
    AND provider_account_id = $2
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	account := &OAuthAccount{}
	err := s.pool.QueryRow(ctx, query, provider, accountID).Scan(
		&account.ID, &account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &account.IDToken,
		&account.AccessTokenExpiresAt, &account.RefreshTokenExpiresAt, &account.Scope,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("oauthAccounts.GetByProviderAndAccountID: %w", ErrNotFound)
		default:
			return nil, fmt.Errorf("oauthAccounts.GetByProviderAndAccountID: %w", err)
		}
	}

	return account, nil
}

func (s *OAuthStore) Create(ctx context.Context, account *OAuthAccount) error {
	query := `
    INSERT INTO oauth_accounts (
      user_id, provider, provider_account_id, access_token, refresh_token, id_token, access_token_expires_at, refresh_token_expires_at, scope
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id, user_id, provider, provider_account_id, access_token, refresh_token, id_token, access_token_expires_at, refresh_token_expires_at, scope, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	err := s.pool.QueryRow(
		ctx, query, account.UserID, account.Provider, account.ProviderAccountID, account.AccessToken,
		account.RefreshToken, account.IDToken, account.AccessTokenExpiresAt,
		account.RefreshTokenExpiresAt, account.Scope,
	).Scan(
		&account.ID, &account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &account.IDToken,
		&account.AccessTokenExpiresAt, &account.RefreshTokenExpiresAt, &account.Scope,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("oauthAccounts.Create: %w", ErrConflict)
		}
		return fmt.Errorf("oauthAccounts.Create: %w", err)
	}

	return nil
}
