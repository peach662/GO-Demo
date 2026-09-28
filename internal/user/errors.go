package user

import "errors"

var ErrUsernameTaken = errors.New("username already taken")