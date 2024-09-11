package comment_api

import (
	middleware "english-ai-full/ecomm-api"
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterCommentRoutes(r *chi.Mux, handler *CommentHandlerController) *chi.Mux {
	tokenMaker := handler.TokenMaker

	r.Route("/comments", func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Comment server is running"))
		})

		// Public routes
		r.Get("/{parentId}", handler.GetComments)
		r.Get("/single/{id}", handler.GetCommentByID)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
			
			r.Post("/", handler.CreateComment)
			r.Put("/", handler.UpdateComment)
			r.Delete("/{id}", handler.DeleteComment)
		})
	})

	return r
}