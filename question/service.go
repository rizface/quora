package question

import (
	"context"
	"errors"

	"github.com/rizface/quora/question/value"
)

type Service struct {
	repo     *Repository
	voteRepo *VoteRepo
}

func NewService(repo *Repository, voteRepo *VoteRepo) *Service {
	return &Service{
		repo:     repo,
		voteRepo: voteRepo,
	}
}

func (s *Service) CreateQuestion(ctx context.Context, q value.QuestionPayload) (value.QuestionEntity, error) {
	var (
		accountId = "f028ac5a-e4c9-442f-bf9a-86c024a79baa" //TODO: update this line using current logged in user
		question  = value.NewQuestionEntity(q, accountId)
	)

	if err := question.Validate(); err != nil {
		return value.QuestionEntity{}, err
	}

	if err := s.repo.Create(ctx, question); err != nil {
		return value.QuestionEntity{}, err
	}

	return question, nil
}

func (s *Service) GetQuestions(ctx context.Context, q value.QuestionQuery) (value.Aggregate, error) {
	if err := value.ValidateQuestionQueery(q); err != nil {
		return value.Aggregate{}, err
	}

	questions, err := s.repo.GetList(ctx, q)
	if err != nil {
		return value.Aggregate{}, err
	}

	totalQuestions, err := s.repo.GetTotalQuestions(ctx)
	if err != nil {
		return value.Aggregate{}, nil
	}

	return value.Aggregate{
		Questions: questions,
		Total:     totalQuestions,
	}, nil
}

func (s *Service) Vote(ctx context.Context, p value.VotePayload) (value.QuestionEntity, error) {
	var (
		voterId = "f028ac5a-e4c9-442f-bf9a-86c024a79baa" //TODO: update this line using current logged in user
		vote    = value.NewVote(p, voterId)
	)

	if err := value.ValidateVote(vote); err != nil {
		return value.QuestionEntity{}, err
	}

	question, err := s.repo.GetOne(ctx, vote.QuestionId)
	if err != nil {
		return value.QuestionEntity{}, err
	}

	oldVote, err := s.voteRepo.GetOldVote(ctx, vote)
	if err != nil && !errors.Is(err, ErrVoteNotFound) {
		return value.QuestionEntity{}, err
	}

	// assume if client spam upvote/downvote button
	if vote.Type == oldVote.Type {
		return question, nil
	}

	question.Vote(vote, oldVote)

	// delete the old vote if the voter had voted the question before
	if oldVote.Type != "" {
		if err := s.voteRepo.DeleteVote(ctx, oldVote); err != nil {
			return value.QuestionEntity{}, err
		}
	}

	if err := s.repo.Vote(ctx, question, vote); err != nil {
		return value.QuestionEntity{}, err
	}

	return question, nil
}
