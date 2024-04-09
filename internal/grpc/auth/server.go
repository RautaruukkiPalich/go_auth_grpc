package auth

import (
	"context"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/utils/validation"
	auth_grpc "github.com/rautaruukkipalich/go_auth_grpc_contract/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, username string, password string) (success bool, err error)
	Login(ctx context.Context, username string, password string, appID int) (token string, err error)
	ChangeUsername(ctx context.Context, username string) (success bool, err error)
	ChangePassword(ctx context.Context, oldPassword string, newPassword string) (success bool, err error)
	Me(ctx context.Context, token string) (user auth_grpc.User, err error)
}

type serverAPI struct {
	auth_grpc.AuthServiceServer
	auth Auth
}

func RegisterServer(gRPC *grpc.Server, auth Auth) {
	auth_grpc.RegisterAuthServiceServer(
		gRPC,
		&serverAPI{auth: auth},
	)
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *auth_grpc.RegisterRequest,
) (*auth_grpc.RegisterResponse, error) {
	
	if err := validateRegister(req.GetPassword(), req.GetUsername()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := s.auth.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_grpc.RegisterResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *auth_grpc.LoginRequest,
) (*auth_grpc.LoginResponse, error) {
	
	if err := validateLogin(req.GetPassword(), req.GetUsername()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: change app id
	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetPassword(), 0)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_grpc.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) ChangePassword(
	ctx context.Context,
	req *auth_grpc.ChangePasswordRequest,
) (*auth_grpc.ChangePasswordResponse, error) {
	if err := validateChangePassword(req.GetToken(), req.GetOldPassword(), req.GetNewPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := s.auth.ChangePassword(ctx, req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_grpc.ChangePasswordResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) ChangeUsername(
	ctx context.Context,
	req *auth_grpc.ChangeUsernameRequest,
) (*auth_grpc.ChangeUsernameResponse, error) {
	if err := validateChangeUsername(req.GetToken(), req.GetUsername()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := s.auth.ChangeUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_grpc.ChangeUsernameResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) Me(
	ctx context.Context,
	req *auth_grpc.MeRequest,
) (*auth_grpc.MeResponse, error) {
	if err := validateMe(req.GetToken()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.auth.Me(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_grpc.MeResponse{
		User: &user,
	}, nil
}

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

func validateChangePassword(token string, oldPassword string, newPassword string) error {
	if err := validation.ValidationToken(token); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidationPassword(oldPassword); err != nil {
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