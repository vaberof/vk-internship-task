package http

import (
	"encoding/json"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

const (
	defaultSearchFilmsLimit  = 100
	defaultSearchFilmsOffset = 0
)

const (
	searchParamFilmTitle = "film-title"
	searchParamActorName = "actor-name"
)

type searchFilmsResponseBody struct {
	Films []*film `json:"films"`
}

//	@Summary		Search films by film`s title or/and actor`s name with optional 'limit' and 'offset' query parameters
//	@Security		BasicAuth
//	@Tags			films
//	@Description	Search films by film title or/and actor name with optional 'limit' and 'offset' query parameters.
//	@Description	If 'film-title' and 'actor-name' are empty, than non-empty list of films with max length = 'limit' will be returned
//	@ID				search-films
//	@Produce		json
//	@Param			film-title	query		string	false	"An optional query parameter 'film-title'"
//	@Param			actor-name	query		string	false	"An optional query parameter 'actor-name'"
//	@Param			limit		query		integer	false	"An optional query parameter 'limit' that limits total number of returned films. By default 'limit' = 100"
//	@Param			offset		query		integer	false	"An optional query parameter 'offset' that indicates how many records should be skipped while listing films. By default 'offset' = 0"
//	@Success		200			{object}	searchFilmsResponseBody
//	@Failure		400			{object}	apiv1.Response
//	@Failure		401			{object}	apiv1.Response
//	@Failure		403			{object}	apiv1.Response
//	@Failure		500			{object}	apiv1.Response
//	@Router			/films/searches [get]
func (h *Handler) SearchFilmsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var limit, offset int
		var err error

		limitStr := request.URL.Query().Get("limit")

		if limitStr == "" {
			limit = defaultSearchFilmsLimit
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			if limit < 0 {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": "'limit' must not be negative"}))

				return
			}
		}

		offsetStr := request.URL.Query().Get("offset")

		if offsetStr == "" {
			offset = defaultSearchFilmsOffset
		} else {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			if offset < 0 {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": "'offset' must not be negative"}))

				return
			}
		}

		filmTitle := request.URL.Query().Get(searchParamFilmTitle)
		actorName := request.URL.Query().Get(searchParamActorName)

		domainFilms, err := h.filmService.SearchByFilters(domain.FilmTitle(filmTitle), domain.ActorName(actorName), limit, offset)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		payload, _ := json.Marshal(&searchFilmsResponseBody{
			Films: buildFilms(domainFilms),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
