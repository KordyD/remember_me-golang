package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/kordyd/remember_me-golang/sso/internal/lib/jwt"
	"github.com/kordyd/remember_me-golang/sso/internal/logger"
	"github.com/kordyd/remember_me-golang/sso/internal/models"
	"github.com/kordyd/remember_me-golang/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type Auth struct {
	userStorage UserStorage
	appStorage  AppStorage
	log         *slog.Logger
	tokenTTL    time.Duration
}

type UserStorage interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	SaveUser(ctx context.Context, email string, passHash []byte) (string, error)
}

type AppStorage interface {
	GetApp(ctx context.Context, id string) (models.App, error)
}

func New(logger *slog.Logger, userStorage UserStorage, appStorage AppStorage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         logger,
		userStorage: userStorage,
		appStorage:  appStorage,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appId string) (string, error) {
	const op = "Auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("email", email), slog.String("appId", appId))
	log.Info("attempting to login")
	user, err := a.userStorage.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid password", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	app, err := a.appStorage.GetApp(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get app", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	token, err := jwt.NewUserJwt(user.Id, user.Email, a.tokenTTL, app.Id, app.Secret)
	if err != nil {
		log.Error("failed to create token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) Register(ctx context.Context, email, password string) (string, error) {
	const op = "Auth.Register"
	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("attempting to register")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	userId, err := a.userStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Warn("user already exists", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		log.Error("failed to save user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
