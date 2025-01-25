package models

import (
	"errors"
	"time"
)

var (
	ErrNoRows = errors.New("there are no rows that satisfy your request")
)

type User struct {
	ID           int
	NickName     string
	Email        string
	PasswordHash string
	SignUp       time.Time
}
type Session struct {
	ID string
}
