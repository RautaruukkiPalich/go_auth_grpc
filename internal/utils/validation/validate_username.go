package validation

import (
	"fmt"
)

var (
	ErrEmptyUsername = fmt.Errorf("empty username")
)

func ValidationUsername(username string) error {
	if username == EmptyString {
		return ErrEmptyUsername
	}

	return nil
}