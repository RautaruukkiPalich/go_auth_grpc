package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/domain/models"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/lib/jwt"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/lib/slerr"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/storage"
	auth_grpc_v1 "github.com/rautaruukkipalich/go_auth_grpc_contract/gen/go/auth"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrGetter  UserGetter
	usrPatcher UserPatcher
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, hashedPass []byte) error
}

type UserGetter interface {
	GetUser(ctx context.Context, username string) (models.User, error)
}

type UserPatcher interface {
	PatchUsername(ctx context.Context, user models.User, username string) error
	PatchPassword(ctx context.Context, user models.User, username string) error
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist = errors.New("user already exists")
)

func New(
	userSaver UserSaver,
	userGetter UserGetter,
	userPatcher UserPatcher,
	appProvider AppProvider,
	log *slog.Logger,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrGetter:   userGetter,
		usrPatcher:  userPatcher,
		appProvider: appProvider,
		log:         log,
		tokenTTL:    tokenTTL,
	}
}

// Register implements auth.Auth.
func (a *Auth) Register(ctx context.Context, username string, password string) (success bool, err error) {
	const op = "auth.Register"
	log := a.log.With(slog.String("op", op))
	log.Info("register user")

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Info("error generating password: ", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err := a.usrSaver.SaveUser(ctx, username, hashedPass); err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			log.Info("user already exists", slerr.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrUserExist)
		}
		log.Error("error save user: ", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

// Login implements auth.Auth.
func (a *Auth) Login(ctx context.Context, username string, password string, appID int) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))
	log.Info("login user")

	user, err := a.usrGetter.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("error get user: ", slerr.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user: ", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPass, []byte(password)); err != nil {
		log.Info("error generating password: ", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Info("error get app: ", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewJWTToken(user, app, a.tokenTTL)
	if err != nil {
		log.Info("error create token: ", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}

// ChangePassword implements auth.Auth.
func (a *Auth) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (success bool, err error) {
	const op = "auth.ChangePassword"
	log := a.log.With(slog.String("op", op))
	log.Info("change password")

	panic("unimplemented")
}

// ChangeUsername implements auth.Auth.
func (a *Auth) ChangeUsername(ctx context.Context, username string) (success bool, err error) {
	const op = "auth.ChangeUsername"
	log := a.log.With(slog.String("op", op))
	log.Info("change username")

	panic("unimplemented")
}

// Me implements auth.Auth.
func (a *Auth) Me(ctx context.Context, token string) (user auth_grpc_v1.User, err error) {
	const op = "auth.Me"
	log := a.log.With(slog.String("op", op))
	log.Info("get me")

	panic("unimplemented")
}
