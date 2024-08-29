package reading_api

import (
	middleware "english-ai-full/ecomm-api"
	"net/http"

	"github.com/go-chi/chi"
)


func RegisterReadingRoutes(r *chi.Mux, handler *ReadingHandlerController) *chi.Mux {
	tokenMaker := handler.TokenMaker
	r.Get("/testreading", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server reading is running"))
	})
	r.Route("/readings", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			
			r.Post("/", handler.CreateReading)
			r.Get("/", handler.ListReadings)
			r.Get("/{id}", handler.FindByID)
			r.Put("/{id}", handler.UpdateReading)
			r.Delete("/{id}", handler.DeleteReading)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			r.Use(middleware.GetAdminMiddlewareFunc(tokenMaker))
			
			r.Get("/page", handler.FindReadingByPage)
		})
	})

	return r
}
