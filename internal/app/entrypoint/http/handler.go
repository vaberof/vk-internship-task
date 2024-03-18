package http

import (
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/middleware/auth"
	"github.com/vaberof/vk-internship-task/internal/domain"
	authservice "github.com/vaberof/vk-internship-task/internal/service/auth"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log/slog"
	"net/http"
)

type Handler struct {
	actorService domain.ActorService
	filmService  domain.FilmService
	authService  authservice.AuthService

	validator *validator.Validate

	logger *slog.Logger
}

func NewHandler(actorService domain.ActorService, filmService domain.FilmService, authService authservice.AuthService, validator *validator.Validate, logsBuilder *logs.Logs) *Handler {
	logger := logsBuilder.WithName("handler")
	return &Handler{
		actorService: actorService,
		filmService:  filmService,
		authService:  authService,
		validator:    validator,
		logger:       logger,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) *http.ServeMux {
	// ====== Actors routes ======

	mux.Handle("POST /api/v1/actors", auth.AuthenticationMiddleware(h.authService,
		auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.CreateActorHandler())))
	mux.Handle("PATCH /api/v1/actors/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.UpdateActorHandler())))
	mux.Handle("DELETE /api/v1/actors/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.DeleteActorHandler())))
	mux.Handle("GET /api/v1/actors", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleUser, user.RoleAdmin}, h.ListActorsHandler())))

	// ====== End of Actors routes ======

	// ====== Films routes ======

	mux.Handle("POST /api/v1/films", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.CreateFilmHandler())))
	mux.Handle("PATCH /api/v1/films/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.UpdateFilmHandler())))
	mux.Handle("DELETE /api/v1/films/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.DeleteFilmHandler())))
	mux.Handle("GET /api/v1/films", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleUser, user.RoleAdmin}, h.ListFilmsHandler())))
	mux.Handle("GET /api/v1/films/searches", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleUser, user.RoleAdmin}, h.SearchFilmsHandler())))

	// ====== End of Films routes ======

	// ====== Swagger route ======

	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	// ====== End of Swagger route ======

	return mux
}
