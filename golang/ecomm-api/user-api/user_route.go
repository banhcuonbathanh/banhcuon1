package User_Api

import (
	middleware "english-ai-full/ecomm-api"
	// middleware "english-ai-full/ecomm-api"

	"net/http"

	"github.com/go-chi/chi"
)



func RegisterRoutesUser(r *chi.Mux, handler *HandlercontrollerUser) *chi.Mux {

	tokenMaker := handler.TokenMaker
	
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server user is running"))
	})
	r.Route("/usersTest", func(r chi.Router) {
		r.Post("/", handler.CreateUsertest)
		r.Post("/login", handler.Logintest)

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetRoleMiddlewareFunc(tokenMaker))
			r.Get("/", handler.ListUserstest)
			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", handler.DeleteUsertest)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			r.Patch("/", handler.Logintest)
			r.Post("/logout", handler.LogoutUsetest)
			r.Get("/email/{email}", handler.FindByEmailtest)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
		r.Route("/user/tokens", func(r chi.Router) {
			r.Post("/user/renew", handler.RenewAccessTokentest)
			r.Post("/user/revoke", handler.RevokeSessiontest)
		})
	})
	return r
}
