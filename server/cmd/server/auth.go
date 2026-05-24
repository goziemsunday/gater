package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chiagxziem/snipper/internal/auth"
	"github.com/chiagxziem/snipper/internal/json"
	"github.com/chiagxziem/snipper/internal/store"
)

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=96"`
}

func (a *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errs, ok := a.validator.ValidateStruct(&payload); !ok {
		json.WriteError(w, http.StatusBadRequest, errs)
		return
	}

	hash, err := auth.HashPassword(payload.Password, nil)
	if err != nil {
		a.logger.Error("internal server error", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user := &store.User{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: &hash,
		Image:        nil,
	}

	err = a.store.Users.Create(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEmailAlreadyExists):
			json.WriteError(w, http.StatusConflict, "email already registered")
		default:
			a.logger.Error("internal server error", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		a.logger.Error("internal server error", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	identifier := "email-verification:" + user.Email
	expiresAt := time.Now().UTC().Add(time.Hour)

	err = a.store.Verifications.Create(
		r.Context(), store.CreateVerificationParams{
			Identifier:  identifier,
			HashedToken: token.Hash,
			ExpiresAt:   expiresAt,
		},
	)
	if err != nil {
		a.logger.Error("internal server error", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	go func() {
		err := a.mailer.SendVerificationEmail(
			context.Background(), []string{user.Email}, user.Name, token.Plaintext,
		)
		if err != nil {
			a.logger.Error("internal server error", "error", err)
		}
	}()

	json.WriteData(w, http.StatusCreated, user)
}

func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) getUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) verifyEmail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) resendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) forgotPwd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) resetPwd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) google(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func (a *application) googleCallback(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}
