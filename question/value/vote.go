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
	AnswerId string
	Type     string // upvote / downvote
}

type Vote struct {
	AnswerId  string
	Type      string
	VoterId   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewVote(p VotePayload, voterId string) Vote {
	return Vote{
		AnswerId: p.AnswerId,
		Type:     p.Type,
		VoterId:  voterId,
	}
}

func ValidateVote(v Vote) error {
	return validation.Errors{
		"type":     validation.Validate(v.Type, validation.Required, validation.In(upvote, downvote)),
		"answerId": validation.Validate(v.AnswerId, validation.Required, is.UUID),
	}.Filter()
}
