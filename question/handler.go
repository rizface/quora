package question

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rizface/quora/question/value"
	"github.com/rizface/quora/stdres"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var payload value.QuestionPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed decode payload",
		})

		return
	}

	question, err := h.svc.CreateQuestion(r.Context(), payload)

	vErr := validation.Errors{}
	if errors.As(err, &vErr) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Data: vErr,
			Info: "validation error",
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
		Data: question,
		Info: "success",
	})
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	// parse query param (limit and skip)
	urlQuery, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed parse url query",
		})

		return
	}

	query, err := value.NewQuestionQuery(urlQuery)
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "invalid query parameter",
		})

		return
	}

	questions, err := h.svc.GetQuestions(r.Context(), query)

	vErr := validation.Errors{}
	if errors.As(err, &vErr) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "validation error",
			Data: vErr,
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
		Info: "success",
		Data: questions,
	})
}
