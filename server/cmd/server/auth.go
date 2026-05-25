package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
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
		json.WriteError(w, http.StatusUnprocessableEntity, errs)
		return
	}

	hash, err := auth.HashPassword(payload.Password, nil)
	if err != nil {
		a.logger.Error("failed to hash password", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
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
		case errors.Is(err, store.ErrConflict):
			json.WriteError(w, http.StatusConflict, "email already registered")
		default:
			a.logger.Error("failed to create user", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		a.logger.Error("failed to generate verification token", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	err = a.store.Verifications.Create(
		r.Context(), store.CreateVerificationParams{
			Identifier:  "email-verification:" + user.Email,
			HashedToken: token.Hash,
			ExpiresAt:   time.Now().UTC().Add(time.Hour),
		},
	)
	if err != nil {
		a.logger.Error("failed to create verification", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	go func() {
		err := a.mailer.SendVerificationEmail(
			context.Background(), []string{user.Email}, user.Name, token.Plaintext,
		)
		if err != nil {
			a.logger.Error("failed to send verification mail", "error", err)
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

type VerifyEmailPayload struct {
	Token string `json:"token" validate:"required,hexadecimal,len=64"`
}

func (a *application) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var payload VerifyEmailPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errs, ok := a.validator.ValidateStruct(&payload); !ok {
		json.WriteError(w, http.StatusUnprocessableEntity, errs)
		return
	}

	hashedToken := auth.HashToken(payload.Token)

	verification, err := a.store.Verifications.Get(r.Context(), hashedToken)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			json.WriteError(w, http.StatusBadRequest, "invalid or expired token")
		default:
			a.logger.Error("failed to get verification", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	// check if token has expired
	if i := time.Now().UTC().Compare(verification.ExpiresAt); i >= 0 {
		// delete token
		if err := a.store.Verifications.Delete(r.Context(), verification.ID); err != nil {
			a.logger.Error("failed to delete verification", "error", err)
		}
		json.WriteError(w, http.StatusBadRequest, "invalid or expired token")
		return
	}

	email, ok := strings.CutPrefix(verification.Identifier, "email-verification:")
	if !ok {
		a.logger.Error("unexpected verification identifier", "identifier", verification.Identifier)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// mark user as verified
	if err := a.store.Users.MarkVerified(r.Context(), email); err != nil {
		a.logger.Error("failed to mark user as verified", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// delete token
	if err := a.store.Verifications.Delete(r.Context(), verification.ID); err != nil {
		a.logger.Error("failed to delete verification", "error", err)
	}

	type returnData struct {
		Status string `json:"status"`
	}
	json.WriteData(w, http.StatusOK, returnData{
		Status: "OK",
	})
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
