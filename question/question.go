package question

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/quora/identifier"
)

type Feature struct {
	handler *Handler
	r       *chi.Mux
}

func NewFeature(r *chi.Mux, db *sql.DB) *Feature {
	var (
		questionRepo = NewRepository(db)
		voteRepo     = NewVoteRepository(db)
		answerRepo   = NewAnswerRepo(db)
		svc          = NewService(questionRepo, voteRepo, answerRepo)
		handler      = NewHandler(svc)
	)

	return &Feature{
		handler: handler,
		r:       r,
	}
}

func (q *Feature) RegisterRoutes() {
	q.r.Group(func(r chi.Router) {
		r.Use(identifier.Identifier)

		r.Route("/questions", func(r chi.Router) {
			r.Post("/", q.handler.CreateQuestion)
			r.Get("/", q.handler.GetQuestion)
		})

		r.Route("/answers", func(r chi.Router) {
			r.Post("/", q.handler.AnswerQuestion)
			r.Patch("/{answerId}/vote", q.handler.Vote)
		})
	})
}
