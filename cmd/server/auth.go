package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/chiagxziem/gater/internal/auth"
	"github.com/chiagxziem/gater/internal/json"
	"github.com/chiagxziem/gater/internal/store"
)

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=96"`
}

func (a *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
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

	type returnData struct {
		Message string      `json:"message"`
		User    *store.User `json:"user"`
	}
	json.WriteData(w, http.StatusCreated, returnData{
		Message: "account created successfully",
		User:    user,
	})
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=96"`
}

func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if errs, ok := a.validator.ValidateStruct(&payload); !ok {
		json.WriteError(w, http.StatusUnprocessableEntity, errs)
		return
	}

	user, err := a.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			json.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		default:
			a.logger.Error("failed to get user by email", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	if user.PasswordHash == nil {
		json.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	matched, err := auth.VerifyPassword(payload.Password, *user.PasswordHash)
	if err != nil {
		a.logger.Error("failed to verify password", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if !matched {
		json.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		a.logger.Error("failed to generate token", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	ua := r.UserAgent()
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		a.logger.Error("failed to get IP address", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	session := &store.Session{
		UserID:    user.ID,
		TokenHash: token.Hash,
		IPAddress: &ip,
		UserAgent: &ua,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 30),
	}

	err = a.store.Sessions.Create(r.Context(), session)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			a.logger.Error("failed to create session: duplicate tokens", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		default:
			a.logger.Error("failed to create session", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	type returnData struct {
		Message string      `json:"message"`
		Token   string      `json:"token"`
		User    *store.User `json:"user"`
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "gater_auth_session",
		Value:    token.Plaintext,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // only over HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30, // 30 days to expiry
	})
	json.WriteData(w, http.StatusOK, returnData{
		Message: "logged in successfully",
		Token:   token.Plaintext,
		User:    user,
	})
}

func (a *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

type VerifyEmailPayload struct {
	Token string `json:"token" validate:"required,hexadecimal,len=64"`
}

func (a *application) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var payload VerifyEmailPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
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
		Message string `json:"message"`
	}
	json.WriteData(w, http.StatusOK, returnData{
		Message: "email verified successfully",
	})
}

type ResendVerificationPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (a *application) resendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var payload ResendVerificationPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if errs, ok := a.validator.ValidateStruct(&payload); !ok {
		json.WriteError(w, http.StatusUnprocessableEntity, errs)
		return
	}

	type returnData struct {
		Message string `json:"message"`
	}
	const successMsg = "if the account with this email exists and is unverified, a verification email has been sent"

	user, err := a.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			json.WriteData(w, http.StatusOK, returnData{Message: successMsg})
		default:
			a.logger.Error("failed to get user by email", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	if user.EmailVerified {
		json.WriteData(w, http.StatusOK, returnData{Message: successMsg})
		return
	}

	identifier := "email-verification:" + user.Email

	latestVerification, err := a.store.Verifications.GetLatest(r.Context(), identifier)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		a.logger.Error("failed to get latest verifications", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	verificationsCount, err := a.store.Verifications.CountSince(r.Context(), identifier, time.Hour)
	if err != nil {
		a.logger.Error("failed to get verifications count", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// only allow resend if there are less than 5 verification tokens
	//  created in an hour and none created in the last 1 minute
	allowResend := latestVerification == nil || (verificationsCount < 5 && time.Now().UTC().Add(-time.Minute).Compare(latestVerification.CreatedAt) == 1)

	if !allowResend {
		json.WriteData(w, http.StatusOK, returnData{Message: successMsg})
		return
	}

	err = a.store.Verifications.DeleteByIdentifier(r.Context(), identifier)
	if err != nil {
		a.logger.Error("failed to delete verification(s)", "error", err)
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

	json.WriteData(w, http.StatusOK, returnData{Message: successMsg})
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
