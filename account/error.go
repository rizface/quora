package account

import "errors"

var (
	ErrEmailIsUsed     = errors.New("email is used")
	ErrUsernameIsUsed  = errors.New("username is used")
	ErrAccountNotFound = errors.New("account not found")
	ErrCredential      = errors.New("wrong email / password")
)
