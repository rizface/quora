package question

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Feature struct {
	handler *Handler
	r       *chi.Mux
}

func NewFeature(r *chi.Mux, db *sql.DB) *Feature {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	return &Feature{
		handler: handler,
		r:       r,
	}
}

func (q *Feature) RegisterRoutes() {
	q.r.Route("/questions", func(r chi.Router) {
		r.Post("/", q.handler.CreateQuestion)
		r.Get("/", q.handler.GetQuestion)
	})
}
