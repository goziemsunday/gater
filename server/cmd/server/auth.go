package main

import (
	"errors"
	"net/http"

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

	newUser, err := a.store.Users.Create(r.Context(), store.CreateUserParams{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: &hash,
		Image:        nil,
	})
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

	// send verification email (this should be a goroutine, no?)
	// 6. Send verification email (goroutine - don't block response)
	// go func() {
	//     a.mailer.SendVerificationEmail(newUser.Email, newUser.ID)
	// }()

	json.WriteData(w, http.StatusCreated, newUser)
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
