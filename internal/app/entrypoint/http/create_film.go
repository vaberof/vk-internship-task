package http

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
	"time"
)

type createFilmRequestBody struct {
	Title       string  `json:"title" validate:"required,min=1,max=150"`
	Description string  `json:"description,omitempty" validate:"max=1000"`
	ReleaseDate string  `json:"release_date" validate:"required"`
	Rating      uint8   `json:"rating" validate:"required,numeric,min=0,max=10"`
	ActorIds    []int64 `json:"actor_ids" validate:"required,gt=0,dive,numeric"`
}

type createFilmResponseBody struct {
	Id          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseDate string   `json:"release_date"`
	Rating      uint8    `json:"rating"`
	Actors      []*actor `json:"actors"`
}

func (h *Handler) CreateFilmHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
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
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageActorNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
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
			Actors:      buildActors(domainFilm.Actors),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
