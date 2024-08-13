package handler

import (
	"english-ai-full/token"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Post("/users", h.CreateUser)
	r.Post("/login", h.LoginUser)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(h.TokenMaker))

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", h.GetUser)
			// Add other user routes (update, delete) as needed
		})

		r.Post("/logout", h.LogoutUser)

		// Add other protected routes as needed
	})

	return r
}

func Start(addr string, router *chi.Mux) error {
	return http.ListenAndServe(addr, router)
}

func AuthMiddleware(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implement JWT token validation logic here
			// If valid, call next.ServeHTTP(w, r)
			// If invalid, return http.Error with appropriate status code
		})
	}
}

// Implement AdminMiddleware if needed
func AdminMiddleware(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implement admin validation logic here
			// If valid admin, call next.ServeHTTP(w, r)
			// If not admin, return http.Error with appropriate status code
		})
	}
}