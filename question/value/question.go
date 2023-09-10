package value

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/rizface/quora/nuller"
)

type QuestionPayload struct {
	SpaceId  nuller.NullString `json:"spaceId"`
	Question string            `json:"question"`
}

type Author struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type QuestionEntity struct {
	Id        string            `json:"id"`
	SpaceId   nuller.NullString `json:"spaceId"`
	AuthorId  string            `json:"authorId"`
	Question  string            `json:"question"`
	Upvote    int               `json:"upvote"`
	Downvote  int               `json:"downvote"`
	Author    Author            `json:"author"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type Aggregate struct {
	Questions []QuestionEntity
	Total     int
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
		"spaceId":  validation.Validate("a53152d7-2d24-42e1-a55f-649e87349ffa", is.UUID),
		"question": validation.Validate(q.Question, validation.Required),
	}.Filter()
}

func (q *QuestionEntity) Vote(vote Vote) {
	if strings.EqualFold(vote.Type, upvote) {
		q.Upvote++
	}

	if strings.EqualFold(vote.Type, downvote) {
		q.Downvote++
	}

	q.UpdatedAt = time.Now()
}
