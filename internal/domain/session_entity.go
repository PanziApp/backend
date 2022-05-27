package domain

import (
	"errors"
	"time"
)

type Session struct {
	Id         EntityId
	CreateTime time.Time
	UserId     EntityId
	Type       TokenType
	Token      Token
	ValidUntil *time.Time
}

const (
	SessionValidUntilFieldName EntityFieldName = "session_valid_until"
)

type Token string

type TokenType string

const (
	GeneralToken           TokenType = "general"
	EmailVerificationToken TokenType = "email-verification"
	ResetPasswordToken     TokenType = "reset-password"
)

var (
	ErrInvalidToken = InternalError{Err: errors.New("invalid token")}
)

func RandomToken() (Token, error) {
	s, err := RandomStringURLSafe(36)
	return Token(s), err
}

func ValidateToken(t string) (Token, error) {
	if len(t) > 100 {
		return "", ErrInvalidToken
	}

	return Token(t), nil
}
