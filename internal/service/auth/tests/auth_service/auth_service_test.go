package auth_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/vaberof/vk-internship-task/internal/service/auth"
	mocks "github.com/vaberof/vk-internship-task/internal/service/auth/mocks"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFinder := mocks.NewMockUserFinder(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	id := int64(1)
	email := "user@example.com"
	password := "asdf1234"
	passwordHash := "$2a$10$JT0HAAksN7kvv6m0TXAvIejUzNOs19uRA7Ae8qIjn5lLa2hP1isNK"
	role := user.RoleUser

	expected := &user.User{
		Id:       id,
		Email:    email,
		Password: passwordHash,
		Role:     role,
	}

	userFinder.EXPECT().FindByEmail(email).Return(expected, nil).Times(1)

	authService := auth.NewAuthService(userFinder, logsBuilder)
	usr, err := authService.AuthenticateUser(email, password)
	require.NoError(t, err)
	require.Equal(t, expected, usr)
}

func TestAuthenticateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userFinder := mocks.NewMockUserFinder(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	authService := auth.NewAuthService(userFinder, logsBuilder)

	type in struct {
		Email    string
		Password string
	}

	testCases := []struct {
		name       string
		in         in
		out        *user.User
		authExpErr error
	}{
		{
			name: "err_invalid_email_or_password",
			in: in{
				Email:    "user@example.com",
				Password: "111",
			},
			authExpErr: auth.ErrInvalidEmailOrPassword,
		},
		{
			name: "err_other",
			in: in{
				Email:    "admin@example.com",
				Password: "aaa",
			},
			authExpErr: fmt.Errorf("failed to authenticate a user: %w", errors.New("failed to hash password")),
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			userFinder.EXPECT().FindByEmail(tCase.in.Email).Return(tCase.out, tCase.authExpErr).AnyTimes()
			usr, err := authService.AuthenticateUser(tCase.in.Email, tCase.in.Password)
			require.Error(t, err)
			require.EqualError(t, tCase.authExpErr, err.Error())
			require.Nil(t, usr)
		})
	}
}
