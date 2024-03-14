package auth

import (
	"github.com/vaberof/vk-internship-task/internal/service/user"
)

type UserFinder interface {
	FindByEmail(email string) (*user.User, error)
}
