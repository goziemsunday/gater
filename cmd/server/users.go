package main

import (
	"net/http"

	"github.com/chiagxziem/gater/internal/json"
	"github.com/chiagxziem/gater/internal/store"
)

const userCtx contextKey = "user"

func (a *application) getUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		a.logger.Error("failed to get user from context")
		json.WriteError(w, http.StatusInternalServerError, "user not found in context")
		return
	}

	type returnData struct {
		User *store.User `json:"user"`
	}
	json.WriteData(w, http.StatusOK, returnData{User: user})
}
