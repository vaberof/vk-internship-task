package user_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	mocks "github.com/vaberof/vk-internship-task/internal/service/user/mocks"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestFindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userStorage := mocks.NewMockUserStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	id := int64(1)
	email := "user@example.com"
	passwordHash := "$2a$10$JT0HAAksN7kvv6m0TXAvIejUzNOs19uRA7Ae8qIjn5lLa2hP1isNK"
	role := user.RoleUser

	expected := &user.User{
		Id:       id,
		Email:    email,
		Password: passwordHash,
		Role:     role,
	}

	userStorage.EXPECT().FindByEmail(email).Return(expected, nil).Times(1)

	userService := user.NewUserService(userStorage, logsBuilder)
	usr, err := userService.FindByEmail(email)
	require.NoError(t, err)
	require.Equal(t, expected, usr)
}

func TestFindByEmailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userStorage := mocks.NewMockUserStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	userService := user.NewUserService(userStorage, logsBuilder)

	type in struct {
		Email string
	}

	testCases := []struct {
		name       string
		in         in
		out        *user.User
		findExpErr error
	}{
		{
			name: "err_user_not_found",
			in: in{
				Email: "user@example.com",
			},
			findExpErr: user.ErrUserNotFound,
		},
		{
			name: "err_other",
			in: in{
				Email: "admin@example.com",
			},
			findExpErr: fmt.Errorf("failed to find a user: %w", errors.New("database is down")),
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			userStorage.EXPECT().FindByEmail(tCase.in.Email).Return(tCase.out, tCase.findExpErr).AnyTimes()
			usr, err := userService.FindByEmail(tCase.in.Email)
			require.Error(t, err)
			require.EqualError(t, tCase.findExpErr, err.Error())
			require.Nil(t, usr)
		})
	}
}
