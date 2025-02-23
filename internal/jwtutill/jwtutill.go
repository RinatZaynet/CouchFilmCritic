package jwtutill

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired  = errors.New("token is expired")
	ErrInvalidToken  = errors.New("token is invalid")
	ErrTokenNotFound = errors.New("token not found")
)

type Claims struct {
	Sub string
	Exp int64
}

type Manager struct {
	privateKey []byte
}

func (manager *Manager) GenJWT(claims *Claims) (token string, err error) {
	const fn = "jwt.GenJWT"

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": claims.Sub,
		"exp": claims.Exp,
	})

	token, err = tokenJWT.SignedString(manager.privateKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return token, nil
}

func (manager *Manager) Parse(token string) (sub string, err error) {
	const fn = "jwt.CheckJWT"
	sub = ""

	if token == "" {
		return sub, ErrTokenNotFound
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: invalid signature method: %v", fn, token.Header["alg"])
		}

		return manager.privateKey, nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return sub, fmt.Errorf("%s: %w", fn, ErrTokenExpired)
		}
		return sub, fmt.Errorf("%s: %w", fn, err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if sub, ok = claims["sub"].(string); ok {
			return sub, nil
		}

		return sub, ErrInvalidToken
	}

	return sub, ErrInvalidToken
}

func NewManager(keyPath string) (*Manager, error) {
	const fn = "jwt.NewManager"

	if keyPath == "" {
		return nil, fmt.Errorf("%s: %w", fn, fmt.Errorf("jwt key length cannot be 0"))
	}

	key := []byte(os.Getenv(keyPath))
	if len(key) == 0 {
		return nil, fmt.Errorf("%s: %w", fn, fmt.Errorf("invalid jwt key"))
	}

	return &Manager{privateKey: key}, nil
}
