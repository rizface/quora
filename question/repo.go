package question

import (
	"context"
	"database/sql"

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
