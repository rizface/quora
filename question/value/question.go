package value

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type QuestionPayload struct {
	SpaceId  *string `json:"spaceId"`
	Question string  `json:"question"`
}

type QuestionEntity struct {
	Id        string    `json:"id"`
	SpaceId   *string   `json:"spaceId"`
	AuthorId  string    `json:"authorId"`
	Question  string    `json:"question"`
	Upvote    int       `json:"upvote"`
	Downvote  int       `json:"downvote"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewQuestionEntity(p QuestionPayload, authorId string) QuestionEntity {
	return QuestionEntity{
		Id:        uuid.NewString(),
		SpaceId:   p.SpaceId,
		AuthorId:  authorId,
		Question:  p.Question,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (q QuestionEntity) Validate() error {
	return validation.Errors{
		"authorId": validation.Validate(q.AuthorId, validation.Required, is.UUID),
		"spaceId":  validation.Validate(q.SpaceId, is.UUID),
		"question": validation.Validate(q.Question, validation.Required),
	}.Filter()
}
