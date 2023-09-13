package question

import (
	"context"
	"database/sql"

	"github.com/rizface/quora/question/value"
)

type (
	AnswerRepo struct {
		db *sql.DB
	}

	CreateAnswerReq struct {
		answer   value.Answer
		question value.QuestionEntity
	}
)

func NewAnswerRepo(db *sql.DB) *AnswerRepo {
	return &AnswerRepo{
		db: db,
	}
}

func (a *AnswerRepo) Create(ctx context.Context, req CreateAnswerReq) (value.Answer, error) {
	var (
		question = req.question
		answer   = req.answer
		command  = `
		INSERT INTO answers (id, question_id, answerer_id, answer, created_at, updated_at) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`
	)

	_, err := a.db.ExecContext(ctx, command, answer.Id, question.Id, answer.AnswererId, answer.Answer, answer.CreatedAt, answer.UpdatedAt)
	if err != nil {
		return value.Answer{}, err
	}

	return answer, nil
}
