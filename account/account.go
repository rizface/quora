package account

import "github.com/go-chi/chi/v5"

type Feature struct {
	Handler *Handler
}

func NewFeature(r *chi.Mux) *Feature {
	return &Feature{
		Handler: NewHandler(r),
	}
}

func (f *Feature) RegisterRoutes() {
	r := f.Handler.r

	r.Get("/", f.Handler.Test)
}
