package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/chiagxziem/gater/internal/auth"
	"github.com/chiagxziem/gater/internal/json"
	"github.com/chiagxziem/gater/internal/store"
)

func (a *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := loggerFromCtx(r.Context())

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			json.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			json.WriteError(w, http.StatusUnauthorized, "malformed authorization header")
			return
		}

		token := parts[1]
		hashedToken := auth.HashToken(token)

		session, err := a.store.Sessions.Get(r.Context(), hashedToken)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				json.WriteError(w, http.StatusUnauthorized, "unauthorized")
			default:
				logger.Error("failed to get session", "error", err, "hashed_token", hashedToken[:8]+"...")
				json.WriteError(w, http.StatusUnauthorized, "unauthorized")
			}
			return
		}

		user, err := a.store.Users.GetByID(r.Context(), session.UserID.String())
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				json.WriteError(w, http.StatusUnauthorized, "unauthorized")
			default:
				logger.Error("failed to get user", "error", err, "session_id", session.ID)
				json.WriteError(w, http.StatusUnauthorized, "unauthorized")
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *application) injectLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		logger := a.logger.With("request_id", reqID)
		ctx := context.WithValue(r.Context(), loggerCtx, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggerFromCtx(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtx).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
