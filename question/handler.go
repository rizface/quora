package question

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
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
			Data: map[string]interface{}{"doc": vErr},
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
		Data: map[string]interface{}{
			"doc": question,
		},
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

	result, err := h.svc.GetQuestions(r.Context(), query)

	vErr := validation.Errors{}
	if errors.As(err, &vErr) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "validation error",
			Data: map[string]interface{}{"doc": vErr},
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
		Data: map[string]interface{}{
			"docs":  result.Questions,
			"total": result.Total,
		},
	})
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	vote := value.VotePayload{
		AnswerId: chi.URLParam(r, "answerId"),
	}

	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed parse payload",
		})

		return
	}

	answer, err := h.svc.Vote(r.Context(), vote)

	vErr := validation.Errors{}
	if errors.As(err, &vErr) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Data: map[string]interface{}{
				"doc": vErr,
			},
			Info: "validation error",
		})

		return
	}

	if errors.Is(err, ErrAnswerNotFound) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusNotFound,
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
		Info: "success",
		Data: map[string]interface{}{
			"doc": answer,
		},
	})
}

func (h *Handler) AnswerQuestion(w http.ResponseWriter, r *http.Request) {
	var (
		payload value.AnswerPayload
	)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed decode answer payload",
		})
	}

	answer, err := h.svc.Answer(r.Context(), payload)

	vErr := validation.Errors{}
	if errors.As(err, &vErr) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Data: map[string]interface{}{"doc": vErr},
		})

		return
	}

	if errors.Is(err, ErrQuestionNotFound) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusNotFound,
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
		Info: "success",
		Data: map[string]interface{}{"doc": answer},
	})
}
