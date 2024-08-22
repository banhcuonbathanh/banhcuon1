package handler

import (
	middleware "english-ai-full/ecomm-api"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)



func RegisterRoutes(r *chi.Mux,handler *handlercontroller) *chi.Mux {

	tokenMaker := handler.TokenMaker
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/", handler.createUser)
		r.Post("/login", handler.login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetAdminMiddlewareFunc(tokenMaker))
			r.Get("/", handler.listUsers)
			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", handler.deleteUser)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			r.Patch("/", handler.login)
			r.Post("/logout", handler.logoutUser)
			r.Get("/email/{email}", handler.FindByEmail)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
		r.Route("/tokens", func(r chi.Router) {
			r.Post("/renew", handler.renewAccessToken)
			r.Post("/revoke", handler.revokeSession)
		})
	})
	return r
}

func Start(addr string, r *chi.Mux) error {
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, r)
}