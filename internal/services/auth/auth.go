package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"
)

var (
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrInvalidCredentials = errors.New("user credentials invalid")
	ErrUserNotFound       = errors.New("user not found")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(
		"op", op,
		"email", email)
	log.Info("register new user", email)
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error pass hash generate", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userId, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("error save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return userId, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int32) (string, error) {
	const op = "auth.Login"
	log := a.log.With(
		"op", op,
		"email", email)

	log.Info("login user", email)

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("error get user", err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		a.log.Warn("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	jwtToken, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("error create token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return jwtToken, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		"op", op,
		"userId", userId,
	)
	log.Info("checking user is admin")
	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err)
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		a.log.Error("error check user is admin", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}
