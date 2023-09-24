package value

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/rizface/quora/identifier"
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
	Author    Author            `json:"author"`
	Answer    Answer            `json:"answer"`
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
		"spaceId":  validation.Validate(q.SpaceId, is.UUID),
		"question": validation.Validate(q.Question, validation.Required),
	}.Filter()
}

func (q QuestionEntity) IsThisTheAuthor(identity identifier.Claim) bool {
	return q.AuthorId == identity.AccountId
}

func (q *QuestionEntity) SyncWithPayload(payload QuestionPayload) {
	q.Question = payload.Question
	q.SpaceId = payload.SpaceId
}
