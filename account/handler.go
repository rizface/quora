package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rizface/quora/account/value"
	"github.com/rizface/quora/stdres"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	r      *chi.Mux
	svc    *Service
	tracer trace.Tracer
}

func NewHandler(r *chi.Mux, svc *Service, tracer trace.Tracer) *Handler {
	return &Handler{
		r:      r,
		svc:    svc,
		tracer: tracer,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "account.Handler.Register")
	defer span.End()

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

	account, err := h.svc.Register(ctx, payload)

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

		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error while create new user: %v", err))

		return
	}

	stdres.Writer(w, stdres.Response{
		Code: http.StatusOK,
		Data: map[string]interface{}{"doc": account},
		Info: "success",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "account.Handler.Login")
	defer span.End()

	var payload value.AccountPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "invalid body request",
		})

		return
	}

	result, err := h.svc.Login(ctx, payload)

	if errors.Is(err, ErrAccountNotFound) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusNotFound,
			Info: err.Error(),
		})

		return
	}

	if errors.Is(err, ErrCredential) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusInternalServerError,
			Info: err.Error(),
		})

		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error while login: %v", err))

		return
	}

	stdres.Writer(w, stdres.Response{
		Code: http.StatusOK,
		Data: result,
		Info: "success",
	})
}
