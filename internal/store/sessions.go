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

type SessionStore struct {
	pool *pgxpool.Pool
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	IPAddress *string   `json:"ip_address"`
	UserAgent *string   `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *SessionStore) Create(ctx context.Context, session *Session) error {
	query := `
    INSERT INTO sessions (user_id, token_hash, ip_address, user_agent, expires_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, user_id, token_hash, ip_address, user_agent, expires_at, created_at,
    updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	err := s.pool.QueryRow(
		ctx, query, session.UserID, session.TokenHash, session.IPAddress,
		session.UserAgent, session.ExpiresAt,
	).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.IPAddress,
		&session.UserAgent, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key"):
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}

func (s *SessionStore) Get(ctx context.Context, hashedToken string) (*Session, error) {
	query := `
    SELECT id, user_id, token_hash, ip_address, user_agent, expires_at, created_at, updated_at
    FROM sessions
    WHERE token_hash = $1
  `

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	session := &Session{}
	err := s.pool.QueryRow(ctx, query, hashedToken).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.IPAddress,
		&session.UserAgent, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return session, nil
}
