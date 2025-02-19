package storage

import (
	"errors"
	"time"
)

var (
	ErrNoRows        = errors.New("no rows found")
	ErrDuplicateData = errors.New("duplicate data")
)

type User struct {
	ID           int
	Nickname     string
	Email        string
	PasswordHash string
	SignUpDate   time.Time
}

type Review struct {
	ID               int
	WorkTitle        string
	Genres           string
	WorkType         string
	Review           string
	Rating           int
	CreateDate       time.Time
	FormatCreateDate string
	Author           string
}
