package validation

import (
	"fmt"
)

var (
	ErrEmptyAppID = fmt.Errorf("empty app id")
)

const (
	ZeroValue = 0
)

func ValidationAppID(appID int32) error {
	if appID == ZeroValue {
		return ErrEmptyAppID
	}

	return nil
}