package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/app/kafka"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/domain/models"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/lib/jwt"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/lib/slerr"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrGetter   UserGetter
	usrPatcher  UserPatcher
	appProvider AppProvider
	tokenTTL    time.Duration
	broker      kafka.Brokerer
}

type UserSaver interface {
	SaveUser(ctx context.Context, email, username string, hashedPass []byte) error
}

type UserGetter interface {
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserPatcher interface {
	PatchUsername(ctx context.Context, user models.User, username string) error
	PatchPassword(ctx context.Context, user models.User, hashed_password []byte) error
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist          = errors.New("user already exists")
)

const (
	ZeroValue = 0
)

func New(
	userSaver UserSaver,
	userGetter UserGetter,
	userPatcher UserPatcher,
	appProvider AppProvider,
	log *slog.Logger,
	tokenTTL time.Duration,
	broker kafka.Brokerer,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrGetter:   userGetter,
		usrPatcher:  userPatcher,
		appProvider: appProvider,
		log:         log,
		tokenTTL:    tokenTTL,
		broker:      broker,
	}
}

// Register implements auth.Auth.
func (a *Auth) Register(ctx context.Context, email, username, password string) (bool, error) {
	const op = "services.auth.Register"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("username", username),
	)
	log.Info("register user")

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Info("error generating password", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err := a.usrSaver.SaveUser(
		ctx,
		strings.ToLower(email),
		username,
		hashedPass,
	); err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			log.Info("user already exists", slerr.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrUserExist)
		}
		log.Error("error save user", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

// Login implements auth.Auth.
func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "services.auth.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.Int("appID", appID),
	)
	log.Info("login user")

	user, err := a.usrGetter.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("failed to get user", slerr.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPass, []byte(password)); err != nil {
		log.Info("failed to check password", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Info("failed to get app", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if app.Secret == "" {
		log.Info("empty secret", slerr.Err(ErrInvalidCredentials))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewJWTToken(user, app, a.tokenTTL)
	if err != nil {
		log.Info("failed to create token", slerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}

// ChangeUsername implements auth.Auth.
func (a *Auth) ChangeUsername(ctx context.Context, token, username string) (bool, error) {
	const op = "services.auth.ChangeUsername"
	log := a.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)
	log.Info("change username")

	appId, err := jwt.GetAppIDFromJWTToken(token)
	if err != nil {
		log.Error("failed to parce app id from token", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		log.Error("failed to get app id", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := jwt.GetSubFromJWTToken(token, app)
	if err != nil {
		log.Error("failed to get user id", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.usrGetter.GetUserByID(ctx, userID)
	if err != nil {
		log.Info("failed to get user from DB", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	err = a.usrPatcher.PatchUsername(ctx, user, username)
	if err != nil {
		log.Info("failed to patch username", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

// ChangePassword implements auth.Auth.
func (a *Auth) ChangePassword(ctx context.Context, token, newPassword string) (bool, error) {
	const op = "services.auth.ChangePassword"
	log := a.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)
	log.Info("change password")

	appId, err := jwt.GetAppIDFromJWTToken(token)
	if err != nil {
		log.Error("failed to parce app id from token", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		log.Error("failed to get app id", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := jwt.GetSubFromJWTToken(token, app)
	if err != nil {
		log.Error("failed to get user id", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.usrGetter.GetUserByID(ctx, userID)
	if err != nil {
		log.Info("failed to get user from DB", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Info("failed to generate password", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	err = a.usrPatcher.PatchPassword(ctx, user, hashedPass)
	if err != nil {
		log.Info("failed to patch password", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

// ResetPassword implements auth.Auth.
func (a *Auth) ResetPassword(ctx context.Context, email string) (bool, error) {
	const op = "services.auth.ResetPassword"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("reset password")

	user, err := a.usrGetter.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("failed to get user", slerr.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	password := generatePassword(email)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Info("failed to generate password", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	err = a.usrPatcher.PatchPassword(ctx, user, hashedPass)
	if err != nil {
		log.Info("failed to patch password", slerr.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	a.broker.AddToQueue(
		kafka.KafkaMessage{
			Topic: "mail",
			Payload: password,
		},
	)

	return true, nil
}

// Me implements auth.Auth.
func (a *Auth) Me(ctx context.Context, token string) (models.User, error) {
	const op = "services.auth.Me"
	log := a.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)
	log.Info("get me")

	var user models.User

	appId, err := jwt.GetAppIDFromJWTToken(token)
	if err != nil {
		log.Error("failed to parce app id from token", slerr.Err(err))
		return user, fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		log.Error("failed to get app id", slerr.Err(err))
		return user, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := jwt.GetSubFromJWTToken(token, app)
	if err != nil {
		log.Error("failed to get user id", slerr.Err(err))
		return user, fmt.Errorf("%s: %w", op, err)
	}

	user, err = a.usrGetter.GetUserByID(ctx, userID)
	if err != nil {
		log.Info("failed to get user from DB", slerr.Err(err))
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
