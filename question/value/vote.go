package value

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	upvote   = "upvote"
	downvote = "downvote"
)

type VotePayload struct {
	QuestionId string
	Type       string // upvote / downvote
}

type Vote struct {
	QuestionId string
	Type       string
	VoterId    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewVote(p VotePayload, voterId string) Vote {
	return Vote{
		QuestionId: p.QuestionId,
		Type:       p.Type,
		VoterId:    voterId,
	}
}

func ValidateVote(v Vote) error {
	return validation.Errors{
		"type":       validation.Validate(v.Type, validation.Required, validation.In(upvote, downvote)),
		"questionId": validation.Validate(v.QuestionId, validation.Required, is.UUID),
	}.Filter()
}
