package value

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type AnswerPayload struct {
	Answer string `json:"answer"`
}

type Answer struct {
	QuestionId string    `json:"questionId"`
	Answer     string    `json:"answer"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewAnswer(questionId string, p AnswerPayload) Answer {
	return Answer{
		QuestionId: questionId,
		Answer:     p.Answer,
	}
}

func ValidateAnswer(a Answer) error {
	return validation.Errors{
		"questionId": validation.Validate(a.QuestionId, validation.Required, is.UUID),
		"answer":     validation.Validate(a.Answer, validation.Required),
	}.Filter()
}
