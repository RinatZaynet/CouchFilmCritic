package auth

import (
	"errors"
)

type Claims struct {
	Sub string
	Exp int64
}

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrInvalidKey    = errors.New("key is invalid")
	ErrTokenExpired  = errors.New("token is expired")
	ErrInvalidToken  = errors.New("token is invalid")
	ErrTokenNotFound = errors.New("token not found")
)
