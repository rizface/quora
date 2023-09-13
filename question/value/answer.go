package value

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type (
	AnswerPayload struct {
		Answer string `json:"answer"`
	}

	Answer struct {
		Id         string    `json:"id"`
		QuestionId string    `json:"questionId"`
		Answer     string    `json:"answer"`
		AnswererId string    `json:"answererId"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	NewAnswerParam struct {
		QuestionId string
		Answer     string
		AnswererId string
	}
)

func NewAnswer(p NewAnswerParam) Answer {
	return Answer{
		Id:         uuid.NewString(),
		QuestionId: p.QuestionId,
		Answer:     p.AnswererId,
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
