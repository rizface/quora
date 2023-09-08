package value

import (
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

type QuestionQuery struct {
	Limit int
	Skip  int
}

func NewQuestionQuery(url url.Values) (QuestionQuery, error) {
	q := QuestionQuery{
		Skip:  0,
		Limit: 20,
	}

	if url.Has("skip") && url.Get("skip") != "" {
		skip, err := strconv.Atoi(url.Get("skip"))
		if err != nil {
			return QuestionQuery{}, err
		}

		q.Skip = skip
	}

	if url.Has("limit") && url.Get("limit") != "" {
		limit, err := strconv.Atoi(url.Get("limit"))
		if err != nil {
			return QuestionQuery{}, err
		}

		q.Limit = limit
	}

	return q, nil
}

func ValidateQuestionQueery(q QuestionQuery) error {
	return validation.Errors{
		"skip":  validation.Validate(q.Skip, validation.Min(0)),
		"limit": validation.Validate(q.Limit, validation.Min(1)),
	}.Filter()
}
