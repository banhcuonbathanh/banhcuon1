package qr_guests

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	// middleware "english-ai-full/ecomm-api"
)

func RegisterGuestRoutes(r *chi.Mux, handler *GuestHandlerController) *chi.Mux {
	// tokenMaker := handler.TokenMaker

	r.Get("/qr/guest/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("guest test is running"))
	})

	r.Route("/qr/guest", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Post("/login", handler.GuestLogin)
		r.Post("/refresh-token", handler.RefreshToken)

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {

			log.Print("golang/quanqr/qr_guests/qr_guest_route.go")
			// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))

			r.Post("/logout", handler.GuestLogout)
			r.Post("/orders", handler.CreateOrders)
			r.Get("/orders/{guestId}", handler.GetOrders)
		})
	})

	return r
}