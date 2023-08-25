package account

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	r *chi.Mux
}

func NewHandler(r *chi.Mux) *Handler {
	return &Handler{
		r: r,
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("root."))
}
