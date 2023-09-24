package value

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type (
	AnswerPayload struct {
		Answer     string `json:"answer"`
		QuestionId string `json:"questionId"`
	}

	Answerer struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}

	Answer struct {
		Id         string    `json:"id"`
		QuestionId string    `json:"questionId,omitempty"`
		Answer     string    `json:"answer"`
		AnswererId string    `json:"answererId,omitempty"`
		Upvote     int       `json:"upvote"`
		Downvote   int       `json:"downvote"`
		Answerer   Answerer  `json:"answerer"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	NewAnswerParam struct {
		AnswerPayload
		AnswererId string
	}
)

func NewAnswer(p NewAnswerParam) Answer {
	return Answer{
		Id:         uuid.NewString(),
		QuestionId: p.QuestionId,
		Answer:     p.Answer,
		AnswererId: p.AnswererId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func ValidateAnswer(a Answer) error {
	return validation.Errors{
		"questionId": validation.Validate(a.QuestionId, validation.Required, is.UUID),
		"answer":     validation.Validate(a.Answer, validation.Required),
	}.Filter()
}

func (q *Answer) Vote(vote Vote, oldVote Vote) {
	if strings.EqualFold(vote.Type, upvote) {
		q.Upvote++

		if strings.EqualFold(oldVote.Type, downvote) {
			q.Downvote--
		}
	}

	if strings.EqualFold(vote.Type, downvote) {
		q.Downvote++

		if strings.EqualFold(oldVote.Type, upvote) {
			q.Upvote--
		}
	}

	q.UpdatedAt = time.Now()
}
