package auth

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
	secretKey   string
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTl time.Duration,
	secretKey string) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		tokenTTL:    tokenTTl,
		secretKey:   secretKey,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op))

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("User not found", "error", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to get user", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("Invalid credentials", "error", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("Successfully logged in", "user", user.Email)

	token, err = jwt.NewToken(user, a.tokenTTL, a.secretKey)
	if err != nil {
		log.Error("Failed to create token", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", "error", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {

			log.Warn("User already exists", "error", err)
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("Failed to save user", "error", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
