package main

import (
	"net/http"

	"github.com/chiagxziem/gater/internal/json"
)

func (a *application) checkHealth(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Status string `json:"status"`
	}
	json.Write(w, http.StatusOK, data{Status: "OK"})
}
