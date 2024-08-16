package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

var r *chi.Mux

func RegisterRoutes(handler *handlercontroller) *chi.Mux {
	r = chi.NewRouter()
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", handler.createUser)
		r.Get("/", handler.listUsers)
		r.Get("/email/{email}", handler.FindByEmail)
		
		// New routes
		r.Delete("/{id}", handler.deleteUser)
		r.Post("/login", handler.login)
		
		// Note: Other user-related routes can be added here as needed
	})

	// Note: Product and Order routes are commented out as they're not implemented in the provided handlercontroller

	return r
}

func Start(addr string) error {
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, r)
}