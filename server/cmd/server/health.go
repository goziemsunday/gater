package main

import "net/http"

func (a *application) checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
