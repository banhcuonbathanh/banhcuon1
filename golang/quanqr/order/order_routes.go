package order_grpc

import (
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterOrderRoutes(r *chi.Mux, handler *OrderHandlerController) *chi.Mux {
	r.Get("/qr/orders-test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Order service is running"))
	})

	r.Route("/qr/orders", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// You can add middleware here if needed
			// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))

			r.Get("/", handler.GetOrders)               // Fetch orders
			r.Post("/", handler.CreateOrders)           // Create new orders

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetOrderDetail)      // Fetch order details by ID
				r.Put("/", handler.UpdateOrder)          // Update an existing order
				// You might want to implement DeleteOrder as well
			})

			r.Post("/pay", handler.PayGuestOrders)      // Pay for guest orders
		})
	})

	return r
}
