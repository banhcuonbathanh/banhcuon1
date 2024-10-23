package route

import (
	middleware "english-ai-full/ecomm-api"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"english-ai-full/ecomm-api/handler"
)



func RegisterRoutes(r *chi.Mux,handler *handler.Handlercontroller) *chi.Mux {

	tokenMaker := handler.TokenMaker
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})
	r.Route("/users", func(r chi.Router) {
		log.Printf("Starting HTTP server /users", )

		r.Post("/", handler.CreateUser)
		r.Post("/login", handler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetRoleMiddlewareFunc(tokenMaker))
			r.Get("/", handler.ListUsers)
			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", handler.DeleteUser)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			r.Patch("/", handler.Login)
			r.Post("/logout", handler.LogoutUser)
			r.Get("/email/{email}", handler.FindByEmail)
		})
	})
	r.Group(func(r chi.Router) {
		// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
		r.Route("/tokens", func(r chi.Router) {
			r.Post("/renew", handler.RenewAccessToken)
			r.Post("/revoke", handler.RevokeSession)
		})
	})
	return r
}

func Start(addr string, r *chi.Mux) error {
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, r)
}