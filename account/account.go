package account

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"
)

type Feature struct {
	Handler *Handler
}

func NewFeature(r *chi.Mux, sql *sql.DB, tracer trace.Tracer) *Feature {
	repo := NewRepository(sql, tracer)
	svc := NewService(repo, tracer)
	handler := NewHandler(r, svc, tracer)

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
