package Python_Api

import (
	// middleware "english-ai-full/ecomm-api"
	"net/http"

	"github.com/go-chi/chi"
)


func RegisterPythonRoutes(r *chi.Mux, handler *PythonHandlerController) *chi.Mux {
	// tokenMaker := handler.TokenMaker
	r.Get("/python-test-server", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server pythontest is running"))
	})
	r.Get("/python-greeter", handler.TestPythonGRPC)
	r.Route("/python", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			
			// r.Post("/", handler.CreateReading)
			// r.Get("/", handler.ListReadings)
			// r.Get("/{id}", handler.FindByID)
			// r.Put("/{id}", handler.UpdateReading)
			// r.Delete("/{id}", handler.DeleteReading)
		})

		r.Group(func(r chi.Router) {
			// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			// r.Use(middleware.GetAdminMiddlewareFunc(tokenMaker))
			
			// r.Get("/page", handler.FindReadingByPage)
		})
	})

	return r
}
