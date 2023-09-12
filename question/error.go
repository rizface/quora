package question

import "errors"

var (
	ErrAuthorNotFound   = errors.New("author not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrVoteNotFound     = errors.New("question not found")
)
