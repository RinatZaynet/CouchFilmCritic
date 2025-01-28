package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/pkg/models"
	"github.com/golang-jwt/jwt"
)

var (
	ErrKeyNotFound    = errors.New("key not found")
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenIsInvalid = errors.New("token is invalid")
	ErrTokenNotFound  = errors.New("token not found")
)

type ManagerJWT struct {
	privateKey []byte
}

func (manager *ManagerJWT) GenTokenJWT(sess *models.Session) (tokenJWT string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sess.Sub,
		"exp": time.Now().Add(time.Hour * 240).Unix(),
	})

	tokenJWT, err = token.SignedString(manager.privateKey)
	if err != nil {
		return "", fmt.Errorf("an error occurred while signing the token in GenTokenJWT. Error: %w", err)
	}

	return tokenJWT, nil
}

func (manager *ManagerJWT) CheckTokenJWT(tokenString string) (sub string, err error) {
	sub = ""
	if tokenString == "" {
		return sub, ErrTokenNotFound
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signature method: %v", token.Header["alg"])
		}
		return manager.privateKey, nil
	})

	if err != nil {
		return sub, fmt.Errorf("an error occurred while parsing the token in CheckTokenJWT. Error: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return sub, ErrTokenExpired
			}
			if sub, ok = claims["sub"].(string); ok {
				return sub, nil
			}
			return sub, ErrTokenIsInvalid
		}
	}
	return sub, ErrTokenIsInvalid
}
func NewClientJWT() (client *ManagerJWT, err error) {
	// Реализовать получаение названия переменной из конф. файла
	key := []byte(os.Getenv("SECRET_KEY"))
	if len(key) == 0 {
		return nil, ErrKeyNotFound
	}
	return &ManagerJWT{privateKey: key}, nil
}
