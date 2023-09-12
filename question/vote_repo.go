package question

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rizface/quora/question/value"
)

type VoteRepo struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) *VoteRepo {
	return &VoteRepo{
		db: db,
	}
}

func (v *VoteRepo) CreateVoter(ctx context.Context, vote value.Vote) (value.Vote, error) {
	return value.Vote{}, nil
}

func (v *VoteRepo) GetOldVote(ctx context.Context, vote value.Vote) (value.Vote, error) {
	var (
		result = value.Vote{}
		query  = `
			SELECT voter_id, question_id, "type", created_at, updated_at FROM votes WHERE voter_id = $1 AND question_id = $2
		`
	)

	err := v.db.
		QueryRowContext(ctx, query, vote.VoterId, vote.QuestionId).
		Scan(&result.VoterId, &result.QuestionId, &result.Type, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return value.Vote{}, ErrVoteNotFound
	}

	if err != nil {
		return value.Vote{}, err
	}

	return result, nil
}

func (v *VoteRepo) DeleteVote(ctx context.Context, vote value.Vote) error {
	command := `
		DELETE FROM votes WHERE voter_id = $1 AND question_id = $2
	`

	if _, err := v.db.ExecContext(ctx, command, vote.VoterId, vote.QuestionId); err != nil {
		return err
	}

	return nil
}
