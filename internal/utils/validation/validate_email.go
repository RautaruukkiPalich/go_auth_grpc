package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrEmptyEmail = fmt.Errorf("empty email")
	ErrInvalidEmail = fmt.Errorf("invalid email")
)

const (
	EmptyString = ""
)

func ValidationEmail(email string) error {
	if email == EmptyString {
		return ErrEmptyEmail
	}

	if strings.Contains(email, " ") {
		return ErrInvalidEmail
	}

	res, err := regexp.MatchString(`^.+@.+\..+$`, email)
	if err != nil {
		return ErrInvalidEmail
	}
	if !res {
		return ErrInvalidEmail
	}

	return nil
}