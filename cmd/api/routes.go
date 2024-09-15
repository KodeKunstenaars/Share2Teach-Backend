package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.authenticate)

	mux.Post("/register", app.registerUser)

	mux.Get("/refresh", app.refreshToken)

	mux.Get("/logout", app.logout)

	mux.Route("/buckets", func(mux chi.Router) {
		// Apply the authRequired middleware to require admin access
		mux.Use(func(next http.Handler) http.Handler {
			return app.authRequired(next, "admin")
		})

		// Route to list buckets (admin-only)
		mux.Get("/", app.listBuckets)
	})

	mux.Route("/upload-document", func(mux chi.Router) {
		// Apply the authRequired middleware to require user access
		mux.Use(func(next http.Handler) http.Handler {
			return app.authRequired(next, "educator", "moderator", "admin")
		})

		// Route to upload a document
		mux.Post("/", app.uploadDocument)
	})

	return mux
}
