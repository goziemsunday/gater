package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/chiagxziem/snipper/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type application struct {
	config *config.Config
	// store
	// cache
	// mailer
	logger *slog.Logger
}

func (a *application) mount() http.Handler {
	r := chi.NewRouter()

	// global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{a.config.CORSAllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// time out reqs after 60 seconds
	r.Use(middleware.Timeout(60 * time.Second))

	// redirect url to short url
	r.Get("/{slug}", a.redirectToShortURL)

	// api routes
	r.Route("/api", func(r chi.Router) {
		// health
		r.Get("/health", a.checkHealth)

		// auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", a.register)
			r.Post("/login", a.login)
			r.Post("/verify-email", a.verifyEmail)
			r.Post("/resend-verification", a.resendVerificationEmail)
			r.Post("/forgot-password", a.forgotPwd)
			r.Post("/reset-password", a.resetPwd)
			r.Get("/google", a.google)
			r.Get("/google/callback", a.googleCallback)

			// protected auth routes
			r.Group(func(r chi.Router) {
				r.Use(a.authed)

				r.Post("/logout", a.logout)
				r.Get("/me", a.getUser)
			})
		})

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(a.authed)

			// urls
			r.Route("/urls", func(r chi.Router) {
				r.Post("/", a.shortenURL)
				r.Get("/", a.listURLs)
				r.Get("/{slug}", a.getURL)
				r.Get("/{slug}/analytics", a.getURLAnalytics)
				r.Patch("/{slug}", a.updateURL)
				r.Delete("/{slug}", a.deleteURL)
			})

			// analytics
			r.Get("/analytics", a.getAnalytics)
		})

	})

	return r
}
