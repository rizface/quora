package account

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rizface/quora/stdres"
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
	stdres.Writer(w, stdres.Response{
		Code:      http.StatusOK,
		Data:      map[string]interface{}{"info": "nice"},
		RequestId: uuid.NewString(),
		Info:      "success",
	})
}
