package models

import (
	"errors"
	"time"
)

var (
	ErrNoRows = errors.New("there are no rows that satisfy your request")
)

type Session struct {
	ID string
}

type User struct {
	ID           int
	NickName     string
	Email        string
	PasswordHash string
	SignUpDate   time.Time
}

type Review struct {
	ID         int
	WorkTitle  string
	Genres     string
	WorkType   string
	Review     string
	Rating     float64
	CreateDate time.Time
	UserID     int
}
