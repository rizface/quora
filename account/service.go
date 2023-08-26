package account

import (
	"context"

	"github.com/rizface/quora/account/value"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, payload value.AccountPayload) (value.AccountEntity, error) {
	account := value.NewAccountEntity(payload)
	if err := account.Validate(); err != nil {
		return account, err
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return account, err
	}

	return account, nil
}
