package question

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
			SELECT id, author_id, question, upvote, downvote, created_at, updated_at FROM questions LIMIT $1 OFFSET $2
		`
	)

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
			&question.Question,
			&question.Upvote,
			&question.Downvote,
			&question.CreatedAt,
			&question.UpdatedAt,
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
			SELECT id, author_id, question, upvote, downvote, created_at, updated_at FROM questions WHERE id = $1
		`
	)

	err := r.db.
		QueryRowContext(ctx, query, questionId).
		Scan(
			&question.Id,
			&question.AuthorId,
			&question.Question,
			&question.Upvote,
			&question.Downvote,
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

func (r *Repository) Vote(ctx context.Context, q value.QuestionEntity, v value.Vote) error {
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
		UPDATE questions SET upvote = $1, downvote = $2, updated_at = $3 WHERE id = $4
	`

	if _, err := tx.ExecContext(ctx, command, q.Upvote, q.Downvote, q.UpdatedAt, q.Id); err != nil {
		return err
	}

	command = `
		INSERT INTO votes (voter_id, question_id, "type") VALUES ($1, $2, $3)
	`

	if _, err := tx.ExecContext(ctx, command, v.VoterId, v.QuestionId, v.Type); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
