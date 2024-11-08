package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(email string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"sub":   "access",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func VerifyAccessToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return false, err
	}
	if sub != "access" {
		return false, errors.New("subject not matching")
	}
	return token.Valid, nil
}

func CreateMFAToken(email string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"sub":   "mfa",
		"id":    id,
		"exp":   time.Now().Add(time.Minute * 10).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_MFA_TOKEN_SECRET")))
}

func VerifyMFAToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_MFA_TOKEN_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return false, err
	}
	if sub != "mfa" {
		return false, errors.New("subject not matching")
	}
	return token.Valid, nil
}

func GetClaimFromToken(tokenString string, secret string, key string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	return int(token.Claims.(jwt.MapClaims)[key].(float64)), nil
}
