package question

import (
	"context"

	"github.com/rizface/quora/question/value"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
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

	// TODO: check if voter already vote the question

	question.Vote(vote)

	if err := s.repo.Vote(ctx, question); err != nil {
		return value.QuestionEntity{}, err
	}

	return question, nil
}
