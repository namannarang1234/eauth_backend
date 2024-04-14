package handlers

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("my secret token")

func CreateToken(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	s, _ := token.SignedString(secret)

	return s
}

func DecryptToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	if m, ok := token.Claims.(jwt.MapClaims); !ok {
		return "", fmt.Errorf("Invalid claims map")
	} else {
		return m["email"].(string), nil
	}
}
