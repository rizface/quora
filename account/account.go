package account

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Feature struct {
	Handler *Handler
}

func NewFeature(r *chi.Mux, sql *sql.DB) *Feature {
	repo := NewRepository(sql)
	svc := NewService(repo)
	handler := NewHandler(r, svc)

	return &Feature{
		Handler: handler,
	}
}

func (f *Feature) RegisterRoutes() {
	r := f.Handler.r

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", f.Handler.Register)
		r.Post("/login", f.Handler.Login)
	})
}
