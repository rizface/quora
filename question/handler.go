package question

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rizface/quora/identifier"
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
	identity, err := identifier.GetFromContext(r.Context())
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

	var payload value.QuestionPayload

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed decode payload",
		})

		return
	}

	question, err := h.svc.CreateQuestion(r.Context(), Input{
		Identity:        *identity,
		QuestionPayload: payload,
	})

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
	identity, err := identifier.GetFromContext(r.Context())
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

	query, err := value.NewQuestionQuery(r.URL.Query())
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "invalid query parameter",
		})

		return
	}

	result, err := h.svc.GetQuestions(r.Context(), Input{
		Identity:      *identity,
		QuestionQuery: query,
	})

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
	identity, err := identifier.GetFromContext(r.Context())
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

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

	answer, err := h.svc.Vote(r.Context(), Input{
		Identity:    *identity,
		VotePayload: vote,
	})

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
	identity, err := identifier.GetFromContext(r.Context())
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

	var (
		payload value.AnswerPayload
	)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed decode answer payload",
		})
	}

	answer, err := h.svc.Answer(r.Context(), Input{
		Identity:      *identity,
		AnswerPayload: payload,
	})

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

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	identity, err := identifier.GetFromContext(r.Context())
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

	input := Input{
		IdQuestion: chi.URLParam(r, "id"),
		Identity:   *identity,
	}

	err = h.svc.DeleteQuestion(r.Context(), input)
	if errors.Is(err, ErrNotTheAuthor) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
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
	})
}

func (h *Handler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		payload    value.QuestionPayload
		idQuestion = chi.URLParam(r, "id")
	)

	identity, err := identifier.GetFromContext(ctx)
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusBadRequest,
			Info: "failed decode payload",
		})

		return
	}

	question, err := h.svc.UpdateQuestion(r.Context(), Input{
		IdQuestion:      idQuestion,
		QuestionPayload: payload,
		Identity:        *identity,
	})

	if errors.Is(err, ErrQuestionNotFound) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusNotFound,
			Info: err.Error(),
		})

		return
	}

	if errors.Is(err, ErrNotTheAuthor) {
		stdres.Writer(w, stdres.Response{
			Code: http.StatusUnauthorized,
			Info: err.Error(),
		})

		return
	}

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
		Data: question,
		Info: "success",
	})
}
