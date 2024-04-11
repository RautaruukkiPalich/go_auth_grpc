package validation

import (
	"fmt"
	"regexp"
)

const (
	minLenPassword = 6
	maxLenPassword = 30
)

var (
	ErrEmptyPassword = fmt.Errorf("empty password")
	ErrInvalidPassword = fmt.Errorf("invalid password")
	ErrInvalidPasswordLength = fmt.Errorf("invalid password. password length must be %d to %d", minLenPassword, maxLenPassword)
)

func ValidationPassword(password string) error {
	if password == EmptyString {
		return ErrEmptyPassword
	}

	pattern := fmt.Sprintf(`^.{%d,%d}$`, minLenPassword, maxLenPassword)
	res, err := regexp.MatchString(pattern, password)
	if err != nil {
		return ErrInvalidPassword
	}
	if !res {
		return ErrInvalidPasswordLength
	}

	return nil
}