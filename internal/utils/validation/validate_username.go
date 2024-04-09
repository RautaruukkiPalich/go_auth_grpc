package validation

import "errors"

var (
	ErrEmptyUsername = errors.New("empty username")
)

func ValidationUsername(username string) error {
	if username == "" {
		return ErrEmptyUsername
	}

	return nil
}