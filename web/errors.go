package web

import "errors"

var (
	NotFound                = errors.New("not found")
	UsernameOrPasswordError = errors.New("username or password error")
	ServerError             = errors.New("server error")
)
