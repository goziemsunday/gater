package main

import (
	"net/http"

	"github.com/chiagxziem/gater/internal/json"
	"github.com/chiagxziem/gater/internal/store"
)

func (a *application) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerFromCtx(ctx)

	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		logger.Error("failed to get user from context")
		json.WriteError(w, http.StatusInternalServerError, "user not found in context")
		return
	}

	type returnData struct {
		User *store.User `json:"user"`
	}
	json.WriteData(w, http.StatusOK, returnData{User: user})
}

func (a *application) becomeOrganizer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerFromCtx(ctx)

	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		logger.Error("failed to get user from context")
		json.WriteError(w, http.StatusInternalServerError, "user not found in context")
		return
	}

	type returnData struct {
		Message string      `json:"message"`
		User    *store.User `json:"user"`
	}

	if user.Role == "organizer" {
		json.WriteData(w, http.StatusOK, returnData{
			Message: "already an organizer",
			User:    user,
		})
		return
	}

	user, err := a.store.Users.BecomeOrganizer(ctx, user.ID.String())
	if err != nil {
		logger.Error("failed to set user as organizer", "error", err, "user_id", user.ID)
		json.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	json.WriteData(w, http.StatusOK, returnData{
		Message: "organizer role activated",
		User:    user,
	})
}
