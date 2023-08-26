package account

import "errors"

var (
	ErrEmailIsUsed    = errors.New("email is used")
	ErrUsernameIsUsed = errors.New("username is used")
)
