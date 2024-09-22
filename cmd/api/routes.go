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

		// Step 1: Route to generate a presigned URL for document upload
		mux.Get("/presigned-url", app.generatePresignedURLForUpload)

		// Step 2: Route to confirm document upload and store metadata
		mux.Post("/confirm", app.uploadDocumentMetadata)
	})

	// mux.Get("/admin/search", app.SearchDocuments)

	mux.Get("/search", app.searchDocuments)

	mux.Route("/admin/search", func(mux chi.Router) {

		mux.Use(func(next http.Handler) http.Handler {
			return app.authRequired(next, "admin", "moderator")
		})

		mux.Get("/", app.searchDocumentsAdminOrModerator)
	})

	mux.Get("/download-document/{id}", app.generatePresignedURLForDownload)

	mux.Get("/faqs", app.FAQs)

	// Route for moderating documents
	mux.Route("/documents/{id}/moderate", func(mux chi.Router) {
		mux.Use(func(next http.Handler) http.Handler {
			return app.authRequired(next, "moderator", "admin")
		})

		mux.Put("/", app.moderateDocument) // Changed from Post to Put
	})

	// Route for rating documents
	mux.Post("/rate-document/{id}", app.rateDocument)

	return mux
}