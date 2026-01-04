package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var tokenSigner *TokenSigner

type TokenSigner struct {
	secret []byte
}

func Init(secret []byte) {
	tokenSigner = &TokenSigner{
		secret: secret,
	}
}

func Get() *TokenSigner {
	return tokenSigner
}

func (s *TokenSigner) GenerateToken(userID string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *TokenSigner) ParseToken(tokenString string) (userID string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	userID, ok = claims["uid"].(string)
	if !ok {
		return "", errors.New("invalid user id")
	}
	return userID, nil
}
