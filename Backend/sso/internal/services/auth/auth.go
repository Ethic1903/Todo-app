package auth

import (
	"AuthService/internal/domain/models"
	"AuthService/internal/lib/jwt"
	"AuthService/internal/storage"
	"context"
	_ "database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type Auth struct {
	log          *slog.Logger
	tokenTTL     time.Duration
	appProvider  AppProvider
	userProvider UserProvider
	userSaver    UserSaver
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserNotExists      = errors.New("user does not exist")
	ErrUserAlrExists      = errors.New("user not found")
)

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{log: log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))
	log.Info("Authorizing user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string) (int64, error) {

	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op))
	log.Info("Registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlrExists) {
			log.Warn("user already exists")
		}
		log.Error("unexpected error", err)
		return 0, fmt.Errorf("%s: %w", op, ErrUserNotExists)
	}
	log.Info("User have been registered")
	return id, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", int(userID)))
	log.Info("checking permissions")

	role, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found")
		}
		return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
	}
	log.Info("checked if user is admin", slog.Bool("is_admin", role))
	return role, nil
}
