package auth

import (
	"errors"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/service/auth"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
)

var (
	ErrMessageUnauthorized        = "errors.middleware.unauthorized"
	ErrMessageForbidden           = "errors.middleware.forbidden"
	ErrMessageInternalServerError = "errors.middleware.internalServerError"
)

func AuthenticationMiddleware(authService auth.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		email, password, hasAuth := request.BasicAuth()
		if !hasAuth {
			views.RenderJSON(writer, http.StatusUnauthorized, apiv1.Error(apiv1.CodeUnauthorized, ErrMessageUnauthorized, apiv1.ErrorDescription{"error": "Missing required 'Authorization' header"}))

			return
		}

		usr, err := authService.AuthenticateUser(email, password)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidEmailOrPassword) {
				views.RenderJSON(writer, http.StatusUnauthorized, apiv1.Error(apiv1.CodeUnauthorized, ErrMessageUnauthorized, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
				views.RenderJSON(writer, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))
			}

			return
		}

		ctxWithUserRole := userRoleToContext(request.Context(), &usr.Role)

		next.ServeHTTP(writer, request.WithContext(ctxWithUserRole))
	})
}

func AuthorizationMiddleware(needRoles []user.UserRole, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userRole := userRoleFromContext(request.Context())
		if userRole == nil {
			views.RenderJSON(writer, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageForbidden, apiv1.ErrorDescription{"error": "Failed to get user role from context"}))

			return
		}

		var hasAccess bool

		for _, needRole := range needRoles {
			if *userRole == needRole {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			views.RenderJSON(writer, http.StatusForbidden, apiv1.Error(apiv1.CodeForbidden, ErrMessageForbidden, apiv1.ErrorDescription{"error": "Access to requested resource has denied"}))

			return
		}

		next.ServeHTTP(writer, request)
	})
}
