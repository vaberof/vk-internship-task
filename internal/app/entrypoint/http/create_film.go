package http

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"log/slog"
	"net/http"
	"time"
)

type createFilmRequestBody struct {
	Title       string  `json:"title" validate:"required,min=1,max=150"`
	Description string  `json:"description,omitempty" validate:"max=1000"`
	ReleaseDate string  `json:"release_date" validate:"required" example:"2006-01-02"`
	Rating      uint8   `json:"rating" validate:"required,numeric,min=0,max=10"`
	ActorIds    []int64 `json:"actor_ids" validate:"required,gt=0,dive,numeric" example:"1,2,3"`
}

type createFilmResponseBody struct {
	Id          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	ReleaseDate string       `json:"release_date"`
	Rating      uint8        `json:"rating"`
	Actors      []*filmActor `json:"actors"`
}

// @Summary		Create a new film
// @Security		BasicAuth
// @Tags			films
// @Description	Create a new film
// @ID				create-film
// @Accept			json
// @Produce		json
// @Param			input	body		createFilmRequestBody	true	"Film object that needs to be created"
// @Success		200		{object}	createFilmResponseBody
// @Failure		400		{object}	apiv1.Response
// @Failure		401		{object}	apiv1.Response
// @Failure		403		{object}	apiv1.Response
// @Failure		404		{object}	apiv1.Response
// @Failure		500		{object}	apiv1.Response
// @Router			/films [post]
func (h *Handler) CreateFilmHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		const handlerName = "CreateFilmHandler"

		log := h.logger.With(slog.String("handlerName", handlerName))

		var createFilmReqBody createFilmRequestBody
		err := json.NewDecoder(request.Body).Decode(&createFilmReqBody)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": "invalid request body"}))

			return
		}

		err = h.validator.Struct(&createFilmReqBody)
		if err != nil {
			errors, ok := err.(validator.ValidationErrors)
			if !ok {
				log.Error("failed to type cast validator.ValidationErrors", "error", err.Error())

				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": "failed to get validation errors"}))
			} else {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": errors.Error()}))
			}

			return
		}

		releaseDate, err := time.Parse(time.DateOnly, createFilmReqBody.ReleaseDate)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		domainFilm, err := h.filmService.Create(
			domain.FilmTitle(createFilmReqBody.Title),
			domain.FilmDescription(createFilmReqBody.Description),
			domain.FilmReleaseDate(releaseDate),
			domain.FilmRating(createFilmReqBody.Rating),
			buildDomainActorIds(createFilmReqBody.ActorIds),
		)
		if err != nil {
			if errors.Is(err, domain.ErrActorsNotFound) {
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageFilmActorsNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
				log.Error("failed to create a film", "error", err.Error())

				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))
			}

			return
		}

		payload, _ := json.Marshal(&createFilmResponseBody{
			Id:          domainFilm.Id.Int64(),
			Title:       domainFilm.Title.String(),
			Description: domainFilm.Description.String(),
			ReleaseDate: domainFilm.ReleaseDate.Time().Format(time.DateOnly),
			Rating:      domainFilm.Rating.Uint8(),
			Actors:      buildFilmActors(domainFilm.Actors),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
