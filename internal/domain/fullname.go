package domain

import "errors"

type Fullname string

var ErrInvalidFullname = ValidationError{Err: errors.New("invalid fullname")}

func ValidateFullname(fullname string) (Fullname, error) {
	if len(fullname) < 5 || 100 < len(fullname) {
		return "", ErrInvalidFullname
	}

	return Fullname(fullname), nil
}
