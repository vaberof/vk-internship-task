package auth

import (
	"context"
	"github.com/vaberof/vk-internship-task/internal/service/user"
)

type userRoleCtxKeyType struct{}

var userRoleCtxKey = &userRoleCtxKeyType{}

func userRoleFromContext(ctx context.Context) *user.UserRole {
	v := ctx.Value(userRoleCtxKey)
	if v == nil {
		return nil
	}

	userRole, ok := v.(*user.UserRole)
	if !ok {
		return nil
	}

	return userRole
}

func userRoleToContext(ctx context.Context, userRole *user.UserRole) context.Context {
	return context.WithValue(ctx, userRoleCtxKey, userRole)
}
