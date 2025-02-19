package random

import (
	"math/rand"
	"time"
)

var chars = []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`)

func RandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}

func RandomSliceByte(size int) []byte {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, size)
	for i := range b {
		b[i] = byte(rnd.Intn(256))
	}

	return b
}
