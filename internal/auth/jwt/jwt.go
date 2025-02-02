package jwt

import (
	"fmt"
	"os"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/golang-jwt/jwt"
)

type ManagerJWT struct {
	privateKey []byte
}

func (manager *ManagerJWT) GenJWT(claims *auth.Claims) (tokenJWT string, err error) {
	const fn = "auth.jwt.GenJWT"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": claims.Sub,
		"exp": claims.Exp,
	})

	tokenJWT, err = token.SignedString(manager.privateKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return tokenJWT, nil
}

func (manager *ManagerJWT) CheckJWT(tokenString string) (sub string, err error) {
	const fn = "auth.jwt.CheckJWT"
	sub = ""

	if tokenString == "" {
		return sub, auth.ErrTokenNotFound
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: invalid signature method: %v", fn, token.Header["alg"])
		}

		return manager.privateKey, nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return sub, fmt.Errorf("%s: %w", fn, auth.ErrTokenExpired)
		}
		return sub, fmt.Errorf("%s: %w", fn, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok = claims["sub"].(string); ok {
			return sub, nil
		}

		return sub, auth.ErrInvalidKey
	}

	return sub, auth.ErrInvalidKey
}

func NewClientJWT(keyPath string) (*ManagerJWT, error) {
	const fn = "auth.jwt.NewClientJWT"

	if keyPath == "" {
		return nil, fmt.Errorf("%s: %w", fn, auth.ErrKeyNotFound)
	}

	key := []byte(os.Getenv(keyPath))
	if len(key) == 0 {
		return nil, fmt.Errorf("%s: %w", fn, auth.ErrInvalidKey)
	}

	return &ManagerJWT{privateKey: key}, nil
}
