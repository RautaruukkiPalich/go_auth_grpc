package validation

import (
	"fmt"
	"regexp"
)

var (
	ErrEmptyToken = fmt.Errorf("empty token")
	ErrInvalidToken = fmt.Errorf("invalid token")
)

func ValidationToken(token string) error {
	if token == EmptyString {
		return ErrEmptyToken
	}

	pattern := `^\w+.\w+.\w+$`
	res, err := regexp.MatchString(pattern, token)
	if err != nil {
		return ErrInvalidToken
	}
	if !res {
		return ErrInvalidToken
	}

	return nil
}