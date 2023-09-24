package value

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	StringIds     []string
	QuestionQuery struct {
		Limit    int
		Skip     int
		SpaceIds StringIds
	}
)

func (s StringIds) ToSqlArray() string {
	if len(s) < 1 {
		return "()"
	}

	ids := ""

	for k, v := range s {
		ids = fmt.Sprintf("%s'%s'", ids, v)
		if k != len(s)-1 {
			ids = fmt.Sprintf("%s, ", ids)
		}
	}

	ids = fmt.Sprintf("(%s)", strings.Trim(ids, " "))

	return ids
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

	if url.Has("space_ids") && len(url.Get("space_ids")) > 0 {
		q.SpaceIds = url["space_ids"]
	}

	return q, nil
}

func ValidateQuestionQueery(q QuestionQuery) error {
	return validation.Errors{
		"skip":  validation.Validate(q.Skip, validation.Min(0)),
		"limit": validation.Validate(q.Limit, validation.Min(1)),
	}.Filter()
}
