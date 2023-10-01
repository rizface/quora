package question

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/quora/identifier"
	"go.opentelemetry.io/otel/trace"
)

type Feature struct {
	handler *Handler
	r       *chi.Mux
}

func NewFeature(r *chi.Mux, db *sql.DB, tracer trace.Tracer) *Feature {
	var (
		questionRepo = NewRepository(db, tracer)
		voteRepo     = NewVoteRepository(db, tracer)
		answerRepo   = NewAnswerRepo(db, tracer)
		svc          = NewService(questionRepo, voteRepo, answerRepo, tracer)
		handler      = NewHandler(svc, tracer)
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
			r.Delete("/{id}", q.handler.DeleteQuestion)
			r.Put("/{id}", q.handler.UpdateQuestion)
		})

		r.Route("/answers", func(r chi.Router) {
			r.Post("/", q.handler.AnswerQuestion)
			// r.Get("/", q.Handler.GetAnswersOfQuestion) -> basically get all answers for specifict question, order by most upvoted
			r.Patch("/{answerId}/vote", q.handler.Vote)
		})
	})
}
