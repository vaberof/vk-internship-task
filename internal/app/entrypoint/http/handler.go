package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/middleware/auth"
	"github.com/vaberof/vk-internship-task/internal/domain"
	authservice "github.com/vaberof/vk-internship-task/internal/service/auth"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	"net/http"
)

type Handler struct {
	actorService domain.ActorService
	filmService  domain.FilmService
	authService  authservice.AuthService

	validator *validator.Validate
}

func NewHandler(actorService domain.ActorService, filmService domain.FilmService, authService authservice.AuthService, validator *validator.Validate) *Handler {
	return &Handler{
		actorService: actorService,
		filmService:  filmService,
		authService:  authService,
		validator:    validator,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) *http.ServeMux {
	// ====== Actors routes ======

	mux.Handle("GET /api/v1/actors", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleUser, user.RoleAdmin}, h.ListActorsHandler())))
	mux.Handle("POST /api/v1/actors", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.CreateActorHandler())))
	mux.Handle("PATCH /api/v1/actors/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.UpdateActorHandler())))
	mux.Handle("DELETE /api/v1/actors/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.DeleteActorHandler())))

	// ====== End of Actors routes ======

	// ====== Films routes ======

	mux.Handle("POST /api/v1/films", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.CreateFilmHandler())))
	mux.Handle("PATCH /api/v1/films/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.UpdateFilmHandler())))
	mux.Handle("DELETE /api/v1/films/{id}", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleAdmin}, h.DeleteFilmHandler())))
	mux.Handle("GET /api/v1/films", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleUser, user.RoleAdmin}, h.ListFilmsHandler())))
	mux.Handle("GET /api/v1/films/searches", auth.AuthenticationMiddleware(h.authService, auth.AuthorizationMiddleware([]user.UserRole{user.RoleUser, user.RoleAdmin}, h.SearchFilmsHandler())))

	// ====== End of Films routes ======

	return mux
}
