package jwt

import "errors"

var (
	ErrJWTDecode = errors.New("invalid token")
)