package user

import (
	"errors"
	"github.com/vaberof/vk-internship-task/internal/infra/storage"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log/slog"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	FindByEmail(email string) (*User, error)
}

type userServiceImpl struct {
	userStorage UserStorage

	logger *slog.Logger
}

func NewUserService(userStorage UserStorage, logsBuilder *logs.Logs) UserService {
	logger := logsBuilder.WithName("domain.service.user")
	return &userServiceImpl{
		userStorage: userStorage,
		logger:      logger,
	}
}

func (u *userServiceImpl) FindByEmail(email string) (*User, error) {
	const operation = "FindByEmail"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", email))

	log.Info("finding a user by email")

	user, err := u.userStorage.FindByEmail(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found by email", "error", err)

			return nil, ErrUserNotFound
		}

		log.Error("failed to find a user", "error", err)

		return nil, err
	}

	log.Info("user has found")

	return user, nil
}
