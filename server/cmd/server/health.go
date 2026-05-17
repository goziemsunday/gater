package main

import "net/http"

func (a *api) checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
