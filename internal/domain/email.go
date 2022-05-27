package domain

import (
	"errors"
	"fmt"
	"regexp"
)

type Email string

var emailRegex = regexp.MustCompile("^[a-zA-Z\\d_.+-]+@[a-zA-Z\\d-]+\\.[a-zA-Z\\d-.]+$")

var ErrInvalidEmail = ValidationError{Err: errors.New("email is not valid")}

func ValidateEmail(email string) (Email, error) {
	if 100 < len(email) ||
		!emailRegex.MatchString(email) {
		return "", ErrInvalidEmail
	}

	return Email(email), nil
}

func ResetPasswordEmailMessage(link string) string {
	return fmt.Sprintf(
		`Hello,<br />
<br />
A reset password request was sent to us. If you didn't send the request, please ignore this email.<br />
<br />
In order to reset your password please click <a href="%s">here</a>.<br />
The link will be valid for the next hour.<br />
<br />
Best Regards,<br />
Fundever Team`,
		link,
	)
}

func EmailVerificationMessage(link string) string {
	return fmt.Sprintf(`Hello,<br />
<br />
In order to verify your email address please click <a href="%s">here</a>.<br />
<br />
If you didn't sign up in Fundever website, please ignore this mail.<br />
<br />
Best Regards,<br />
Fundever Team`,
		link,
	)
}
