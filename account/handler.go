package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rizface/quora/account/value"
	"github.com/rizface/quora/stdres"
)

type Handler struct {
	r   *chi.Mux
	svc *Service
}

func NewHandler(r *chi.Mux, svc *Service) *Handler {
	return &Handler{
		r:   r,
		svc: svc,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		payload value.AccountPayload
		err     error
	)

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "invalid request body",
		})

		return
	}

	account, err := h.svc.Register(r.Context(), payload)

	if errors.As(err, &validation.Errors{}) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Data: err,
			Info: "validation error",
		})

		return
	}

	if errors.Is(err, ErrEmailIsUsed) || errors.Is(err, ErrUsernameIsUsed) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusConflict,
			Data: account,
			Info: err.Error(),
		})

		return
	}

	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusInternalServerError,
			Info: err.Error(),
		})

		return
	}

	stdres.Writer(w, stdres.Response{
		Code: http.StatusOK,
		Data: map[string]interface{}{"doc": account},
		Info: "success",
	})
}
