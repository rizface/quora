package question

import (
	"context"
	"database/sql"
	"errors"
	"log"

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

func (a *AnswerRepo) GetOne(ctx context.Context, answerId string) (value.Answer, error) {
	var (
		answer = value.Answer{}
		query  = `
			SELECT id, question_id, answerer_id, answer, upvote, downvote, created_at, updated_at FROM answers WHERE id = $1
		`
	)

	err := a.db.QueryRowContext(ctx, query, answerId).
		Scan(
			&answer.Id,
			&answer.QuestionId,
			&answer.AnswererId,
			&answer.Answer,
			&answer.Upvote,
			&answer.Downvote,
			&answer.CreatedAt,
			&answer.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return value.Answer{}, ErrAnswerNotFound
	}

	if err != nil {
		return value.Answer{}, err
	}

	return answer, nil
}

func (r *AnswerRepo) Vote(ctx context.Context, q value.Answer, v value.Vote) error {
	tx, err := r.db.Begin()

	defer func(err error) {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
			}

			return
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)

			return
		}
	}(err)

	if err != nil {
		return err
	}

	command := `
		UPDATE answers SET upvote = $1, downvote = $2, updated_at = $3 WHERE id = $4
	`

	if _, err := tx.ExecContext(ctx, command, q.Upvote, q.Downvote, q.UpdatedAt, q.Id); err != nil {
		return err
	}

	command = `
		INSERT INTO votes (voter_id, answer_id, "type") VALUES ($1, $2, $3)
	`

	if _, err := tx.ExecContext(ctx, command, v.VoterId, v.AnswerId, v.Type); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
