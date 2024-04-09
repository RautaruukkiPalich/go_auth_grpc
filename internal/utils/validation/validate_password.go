package validation

import "errors"

var (
	ErrEmptyPassword = errors.New("empty password")
)

func ValidationPassword(pwrd string) error {
	if pwrd == "" {
		return ErrEmptyPassword
	}

	return nil
}