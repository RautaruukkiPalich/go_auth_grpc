package auth

import (
	"context"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/domain/models"
	auth_grpc "github.com/rautaruukkipalich/go_auth_grpc_contract/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, username, password string) (success bool, err error)
	Login(ctx context.Context, username, password string, appID int) (token string, err error)
	ChangeUsername(ctx context.Context, token, username string) (success bool, err error)
	ChangePassword(ctx context.Context, token, newPassword string) (success bool, err error)
	Me(ctx context.Context, token string) (user models.User, err error)
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

	// TODO: change app id get from req
	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetPassword(), 1)
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
	if err := validateChangePassword(req.GetToken(), req.GetNewPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := s.auth.ChangePassword(ctx, req.GetToken(), req.GetNewPassword())
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

	success, err := s.auth.ChangeUsername(ctx, req.GetToken(), req.GetUsername())
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
		User: &auth_grpc.User{
			Id: user.ID,
			Username: user.Username,
		},
	}, nil
}
