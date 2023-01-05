package helper

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

func ParseTokenToClaims(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("error when asserting")
	}

	return claims, nil
}
