package question

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rizface/quora/question/value"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, q value.QuestionEntity) error {
	query := `
		INSERT INTO questions (id, author_id, space_id, question) VALUES($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, q.Id, q.AuthorId, q.SpaceId, q.Question)

	return err
}

func (r *Repository) GetList(ctx context.Context, q value.QuestionQuery) ([]value.QuestionEntity, error) {
	var (
		questions = []value.QuestionEntity{}
		query     = `
			SELECT 
			q.id, q.author_id, q.space_id, q.question, q.created_at, q.updated_at, a.id, a.answer, 
			a.upvote, a.downvote, a.created_At, a.updated_at,
			ac.id, ac.username 
			FROM questions q
			INNER JOIN answers a ON a.question_id = q.id
			INNER JOIN accounts ac ON ac.id = a.answerer_id
		`
	)

	if len(q.SpaceIds) > 0 {
		query = fmt.Sprintf("%s %s %s", query, "WHERE space_id IN", q.SpaceIds.ToSqlArray())
	}

	query = fmt.Sprintf("%s %s %s", query, "ORDER BY a.updated_at DESC", "LIMIT $1 OFFSET $2")

	rows, err := r.db.QueryContext(ctx, query, q.Limit, q.Skip)
	if err != nil {
		return []value.QuestionEntity{}, err
	}

	var getAuthor = func(authorId string) (value.Author, error) {
		var (
			author = value.Author{}
			query  = `SELECT id, username FROM accounts WHERE id = $1`
		)

		err := r.db.
			QueryRowContext(ctx, query, authorId).
			Scan(&author.Id, &author.Username)
		if errors.Is(err, sql.ErrNoRows) {
			return author, ErrAuthorNotFound
		}

		return author, err
	}

	for rows.Next() {
		question := value.QuestionEntity{}

		err := rows.Scan(
			&question.Id,
			&question.AuthorId,
			&question.SpaceId,
			&question.Question,
			&question.CreatedAt,
			&question.UpdatedAt,
			&question.Answer.Id,
			&question.Answer.Answer,
			&question.Answer.Upvote,
			&question.Answer.Downvote,
			&question.Answer.CreatedAt,
			&question.Answer.UpdatedAt,
			&question.Answer.Answerer.Id,
			&question.Answer.Answerer.Username,
		)
		if err != nil {
			return []value.QuestionEntity{}, err
		}

		question.Author, err = getAuthor(question.AuthorId)
		if err != nil {
			return []value.QuestionEntity{}, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (r *Repository) GetTotalQuestions(ctx context.Context) (int, error) {
	var (
		total int
		query = `
			SELECT COUNT(*) FROM questions
		`
	)

	if err := r.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		return total, err
	}

	return total, nil
}

func (r *Repository) GetOne(ctx context.Context, questionId string) (value.QuestionEntity, error) {
	var (
		question value.QuestionEntity
		query    = `
			SELECT id, author_id, question, created_at, updated_at FROM questions WHERE id = $1
		`
	)

	err := r.db.
		QueryRowContext(ctx, query, questionId).
		Scan(
			&question.Id,
			&question.AuthorId,
			&question.Question,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return value.QuestionEntity{}, ErrQuestionNotFound
	}

	if err != nil {
		return value.QuestionEntity{}, err
	}

	return question, nil
}

func (r *Repository) DeleteQuestion(ctx context.Context, question value.QuestionEntity) error {
	command := `
		DELETE FROM questions WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, command, question.Id)

	return err
}

func (r *Repository) UpdateQuestion(ctx context.Context, question value.QuestionEntity) error {
	command := `
		UPDATE questions SET question = $1, space_id = $2 WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, command, question.Question, question.SpaceId, question.Id)

	return err
}
