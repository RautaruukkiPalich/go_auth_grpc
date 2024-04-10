package auth

import (
	"github.com/rautaruukkipalich/go_auth_grpc/internal/utils/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateRegister(username string, password string) error {
	if err := validation.ValidationUsername(username); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidationPassword(password); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func validateLogin(username string, password string) error {
	if err := validation.ValidationUsername(username); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidationPassword(password); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func validateChangePassword(token string, newPassword string) error {
	if err := validation.ValidationToken(token); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidationPassword(newPassword); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func validateChangeUsername(token string, username string) error {
	if err := validation.ValidationUsername(username); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidationToken(token); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func validateMe(token string) error {
	if err := validation.ValidationToken(token); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}