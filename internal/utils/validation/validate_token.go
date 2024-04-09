package validation

import "errors"

var (
	ErrEmptyToken = errors.New("empty token")
)

func ValidationToken(token string) error {
	if token == "" {
		return ErrEmptyToken
	}

	return nil
}