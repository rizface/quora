package question

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rizface/quora/question/value"
	"go.opentelemetry.io/otel/trace"
)

type VoteRepo struct {
	db     *sql.DB
	tracer trace.Tracer
}

func NewVoteRepository(db *sql.DB, tracer trace.Tracer) *VoteRepo {
	return &VoteRepo{
		db:     db,
		tracer: tracer,
	}
}

func (v *VoteRepo) GetOldVote(ctx context.Context, vote value.Vote) (value.Vote, error) {
	ctx, span := v.tracer.Start(ctx, "question.VoteRepo.GetOldVote")
	defer span.End()

	var (
		result = value.Vote{}
		query  = `
			SELECT voter_id, answer_id, "type", created_at, updated_at FROM votes WHERE voter_id = $1 AND answer_id = $2
		`
	)

	err := v.db.
		QueryRowContext(ctx, query, vote.VoterId, vote.AnswerId).
		Scan(&result.VoterId, &result.AnswerId, &result.Type, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return value.Vote{}, ErrVoteNotFound
	}

	if err != nil {
		return value.Vote{}, err
	}

	return result, nil
}

func (v *VoteRepo) DeleteVote(ctx context.Context, vote value.Vote) error {
	ctx, span := v.tracer.Start(ctx, "question.VoteRepo.DeleteVote")
	defer span.End()

	command := `
		DELETE FROM votes WHERE voter_id = $1 AND answer_id = $2
	`

	if _, err := v.db.ExecContext(ctx, command, vote.VoterId, vote.AnswerId); err != nil {
		return err
	}

	return nil
}
