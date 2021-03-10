package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) initializeRouter() {
	s.router = chi.NewRouter()

	s.router.Use(
		middleware.Logger,
		middleware.RedirectSlashes,
	)

	s.router.Route("/todos", func(r chi.Router) {
		r.Post("/", s.controller.CreateToDo())
		r.Get("/", s.controller.GetToDos())

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.controller.GetToDo())
			r.Put("/", s.controller.UpdateToDo())
			r.Delete("/", s.controller.DeleteToDo())
		})
	})
}
