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

func (s *Service) CreateQuestion(ctx context.Context, q value.QuestionPayload) error {
	var (
		accountId = "f028ac5a-e4c9-442f-bf9a-86c024a79baa" //TODO: update this line using current logged in user
		question  = value.NewQuestionEntity(q, accountId)
	)

	if err := question.Validate(); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, question); err != nil {
		return err
	}

	return nil
}
