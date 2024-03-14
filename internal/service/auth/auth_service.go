package auth

import (
	"errors"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"github.com/vaberof/vk-internship-task/pkg/xpassword"
	"log/slog"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

type AuthService interface {
	AuthenticateUser(email, password string) error
}

type authServiceImpl struct {
	userFinder UserFinder

	logger *slog.Logger
}

func NewAuthService(userFinder UserFinder, logsBuilder *logs.Logs) AuthService {
	logger := logsBuilder.WithName("domain.service.auth")
	return &authServiceImpl{
		userFinder: userFinder,
		logger:     logger,
	}
}

func (a *authServiceImpl) AuthenticateUser(email, password string) error {
	const operation = "AuthenticateUser"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email))

	log.Info("authenticating a user")

	usr, err := a.userFinder.FindByEmail(email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			log.Error("failed to authenticate a user", "error", user.ErrUserNotFound)
			return ErrInvalidEmailOrPassword
		}

		log.Error("failed to authenticate a user", "error", err)
		return err
	}

	if err = xpassword.Check(password, usr.Password); err != nil {
		log.Error("failed to authenticate a user", "error", err)
		return ErrInvalidEmailOrPassword
	}

	log.Info("user has authenticated")

	return nil
}
