package user

import "errors"

var (
	ErrLoginTaken = errors.New("login_taken")
)
