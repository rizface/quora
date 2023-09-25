package account

import (
	"context"

	"github.com/rizface/quora/account/value"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	repo   *Repository
}

func NewService(repo *Repository, tracer trace.Tracer) *Service {
	return &Service{
		repo:   repo,
		tracer: tracer,
	}
}

func (s *Service) Register(ctx context.Context, payload value.AccountPayload) (value.AccountEntity, error) {
	ctx, span := s.tracer.Start(ctx, "account.Service.Register")
	defer span.End()

	account := value.NewAccountEntity(payload)
	if err := account.Validate(); err != nil {
		return account, err
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return account, err
	}

	return account, nil
}
func (s *Service) Login(ctx context.Context, payload value.AccountPayload) (value.Authenticated, error) {
	ctx, span := s.tracer.Start(ctx, "account.Service.Login")
	defer span.End()

	account := value.NewAccountEntity(payload)

	account, err := s.repo.FindByEmail(ctx, account)
	if err != nil {
		return value.Authenticated{}, err
	}

	if !account.VerifyPassword(payload.Password) {
		return value.Authenticated{}, ErrCredential
	}

	authenticated, err := value.NewAuthenticated(account)

	return authenticated, err
}
