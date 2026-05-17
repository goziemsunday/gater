package main

import "net/http"

func (a *api) authed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure user is authenticated and pass to context
		// ctx := context.WithValue(r.Context(), "user", "123")

		// next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}
