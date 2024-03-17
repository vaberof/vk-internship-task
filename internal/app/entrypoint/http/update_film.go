package http

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
	"time"
)

type updateFilmRequestBody struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=150"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	ReleaseDate *string `json:"release_date,omitempty"`
	Rating      *uint8  `json:"rating,omitempty" validate:"omitempty,numeric,min=0,max=10"`
}

type updateFilmResponseBody struct {
	Id          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	ReleaseDate string       `json:"release_date"`
	Rating      uint8        `json:"rating"`
	Actors      []*filmActor `json:"actors"`
}

//	@Summary		Update fully or partially a film by path parameter 'id'
//	@Security		BasicAuth
//	@Tags			films
//	@Description	Update fully or partially a film by path parameter 'id'
//	@ID				update-film
//	@Accept			json
//	@Produce		json
//	@Param			id		path		integer					true	"Films`s id that needs to be updated"
//	@Param			input	body		updateFilmRequestBody	true	"Film object with values that will be updated"
//	@Success		200		{object}	updateFilmResponseBody
//	@Failure		400		{object}	apiv1.Response
//	@Failure		401		{object}	apiv1.Response
//	@Failure		403		{object}	apiv1.Response
//	@Failure		404		{object}	apiv1.Response
//	@Failure		500		{object}	apiv1.Response
//	@Router			/films/{id} [patch]
func (h *Handler) UpdateFilmHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var updateFilmReqBody updateFilmRequestBody
		err := json.NewDecoder(request.Body).Decode(&updateFilmReqBody)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": "invalid request body"}))

			return
		}

		err = h.validator.Struct(&updateFilmReqBody)
		if err != nil {
			errors, ok := err.(validator.ValidationErrors)
			if !ok {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": "failed to get validation errors"}))
			} else {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": errors.Error()}))
			}

			return
		}

		var domainReleaseDate *domain.FilmReleaseDate
		if updateFilmReqBody.ReleaseDate != nil {
			parsedReleaseDate, err := time.Parse(time.DateOnly, *updateFilmReqBody.ReleaseDate)
			if err != nil {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			convDomainBirthdate := domain.FilmReleaseDate(parsedReleaseDate)
			domainReleaseDate = &convDomainBirthdate
		}

		filmIdPathParam := request.PathValue("id")
		filmId, err := strconv.Atoi(filmIdPathParam)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		var domainFilmTitle *domain.FilmTitle
		if updateFilmReqBody.Title != nil {
			convDomainFilmTitle := domain.FilmTitle(*updateFilmReqBody.Title)
			domainFilmTitle = &convDomainFilmTitle
		}

		var domainFilmDescription *domain.FilmDescription
		if updateFilmReqBody.Description != nil {
			convDomainFilmDescription := domain.FilmDescription(*updateFilmReqBody.Description)
			domainFilmDescription = &convDomainFilmDescription
		}

		var domainFilmRating *domain.FilmRating
		if updateFilmReqBody.Rating != nil {
			convDomainFilmRating := domain.FilmRating(*updateFilmReqBody.Rating)
			domainFilmRating = &convDomainFilmRating
		}

		domainFilm, err := h.filmService.Update(
			domain.FilmId(filmId),
			domainFilmTitle,
			domainFilmDescription,
			domainReleaseDate,
			domainFilmRating,
		)
		if err != nil {
			if errors.Is(err, domain.ErrFilmNotFound) {
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageFilmNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))
			}

			return
		}

		payload, _ := json.Marshal(&updateFilmResponseBody{
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
