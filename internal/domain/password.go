package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	Password       []byte
	HashedPassword []byte
)

var ErrShortPassword = ValidationError{Err: errors.New("password should be at least 8 characters")}

func ValidatePassword(password string) (Password, error) {
	if len(password) < 8 || 100 < len(password) {
		return Password{}, ErrShortPassword
	}

	return Password(password), nil
}

func HashPassword(password Password) (HashedPassword, error) {
	if len(password) < 8 || 100 < len(password) {
		return HashedPassword{}, ErrShortPassword
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, InternalError{Err: err}
	}

	return h, nil
}

var ErrInvalidPassword = ValidationError{Err: errors.New("invalid password")}

func (p HashedPassword) Match(password Password) error {
	err := bcrypt.CompareHashAndPassword(p, password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidPassword
	} else if err != nil {
		return InternalError{Err: err}
	}

	return nil
}
