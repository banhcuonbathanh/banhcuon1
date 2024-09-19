package Python_Api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterPythonIeltsRoutes(r *chi.Mux, handler *PythonIeltsHandlerController) *chi.Mux {
	r.Get("/python-test-server-ielts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server pythontest is running"))
	})

	r.Route("/ielts", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// You can add middleware here if needed
			// r.Use(middleware.SomeMiddleware)

			// Add the TestPythonGRPC handler
			r.Get("/evaluate", handler.TestPythonGRPC)

			// You can add more routes here as needed
			// For example:
			// r.Post("/submit", handler.SubmitIELTS)
			// r.Get("/results", handler.GetIELTSResults)
		})
	})

	return r
}