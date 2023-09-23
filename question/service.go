package question

import (
	"context"
	"errors"

	"github.com/rizface/quora/identifier"
	"github.com/rizface/quora/question/value"
)

type (
	Service struct {
		repo       *Repository
		voteRepo   *VoteRepo
		answerRepo *AnswerRepo
	}

	AnwerQuestionRequest struct {
		value.AnswerPayload
	}

	Input struct {
		IdQuestion      string
		Identity        identifier.Claim
		QuestionPayload value.QuestionPayload
		QuestionQuery   value.QuestionQuery
		VotePayload     value.VotePayload
		AnswerPayload   value.AnswerPayload
	}
)

func NewService(repo *Repository, voteRepo *VoteRepo, answerRepo *AnswerRepo) *Service {
	return &Service{
		repo:       repo,
		voteRepo:   voteRepo,
		answerRepo: answerRepo,
	}
}

func (s *Service) CreateQuestion(ctx context.Context, input Input) (value.QuestionEntity, error) {
	var (
		accountId = input.Identity.AccountId
		question  = value.NewQuestionEntity(input.QuestionPayload, accountId)
	)

	if err := question.Validate(); err != nil {
		return value.QuestionEntity{}, err
	}

	if err := s.repo.Create(ctx, question); err != nil {
		return value.QuestionEntity{}, err
	}

	return question, nil
}

func (s *Service) GetQuestions(ctx context.Context, input Input) (value.Aggregate, error) {
	if err := value.ValidateQuestionQueery(input.QuestionQuery); err != nil {
		return value.Aggregate{}, err
	}

	questions, err := s.repo.GetList(ctx, input.QuestionQuery)
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

func (s *Service) Vote(ctx context.Context, input Input) (value.Answer, error) {
	var (
		voterId = input.Identity.AccountId
		vote    = value.NewVote(input.VotePayload, voterId)
	)

	if err := value.ValidateVote(vote); err != nil {
		return value.Answer{}, err
	}

	answer, err := s.answerRepo.GetOne(ctx, vote.AnswerId)
	if err != nil {
		return value.Answer{}, err
	}

	oldVote, err := s.voteRepo.GetOldVote(ctx, vote)
	if err != nil && !errors.Is(err, ErrVoteNotFound) {
		return value.Answer{}, err
	}

	// assume if client spam upvote/downvote button
	if vote.Type == oldVote.Type {
		return answer, nil
	}

	answer.Vote(vote, oldVote)

	// delete the old vote if the voter had voted the question before
	if oldVote.Type != "" {
		if err := s.voteRepo.DeleteVote(ctx, oldVote); err != nil {
			return value.Answer{}, err
		}
	}

	if err := s.answerRepo.Vote(ctx, answer, vote); err != nil {
		return value.Answer{}, err
	}

	return answer, nil
}

func (s *Service) Answer(ctx context.Context, input Input) (value.Answer, error) {
	var (
		answererId = input.Identity.AccountId
		answer     = value.NewAnswer(value.NewAnswerParam{
			AnswerPayload: input.AnswerPayload,
			AnswererId:    answererId,
		})
	)

	if err := value.ValidateAnswer(answer); err != nil {
		return value.Answer{}, err
	}

	question, err := s.repo.GetOne(ctx, answer.QuestionId)
	if err != nil {
		return value.Answer{}, err
	}

	answer, err = s.answerRepo.Create(ctx, CreateAnswerReq{
		answer:   answer,
		question: question,
	})
	if err != nil {
		return value.Answer{}, err
	}

	return answer, nil
}

func (s *Service) DeleteQuestion(ctx context.Context, input Input) error {
	question, err := s.repo.GetOne(ctx, input.IdQuestion)
	if err != nil {
		return err
	}

	if !question.IsThisTheAuthor(input.Identity) {
		return ErrNotTheAuthor
	}

	if err = s.repo.DeleteQuestion(ctx, question); err != nil {
		return err
	}

	return nil
}
