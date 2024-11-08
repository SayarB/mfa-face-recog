package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mfa-face-recog/pkg/auth/config"
)

type SessionToken struct {
	ID    string `json:"session_id"`
	Token string `json:"token"`
}

func CreateMFASession(id int) (*SessionToken, error) {
	var sessionID int
	err := config.DB.QueryRow(`INSERT INTO mfa_sessions (user_id) VALUES ($1) RETURNING id`, id).Scan(&sessionID)
	if err != nil {
		return nil, err
	}

	sessionToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      sessionID,
		"user_id": id,
		"sub":     "mfa-session",
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	})
	tok, err := sessionToken.SignedString([]byte(os.Getenv("JWT_MFA_SESSION_SECRET")))
	if err != nil {
		return nil, err
	}
	return &SessionToken{ID: strconv.Itoa(sessionID), Token: tok}, nil
}

func VerifyMFASession(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_MFA_SESSION_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return false, err
	}
	if sub != "mfa-session" {
		return false, errors.New("subject not matching")
	}
	return token.Valid, nil
}
