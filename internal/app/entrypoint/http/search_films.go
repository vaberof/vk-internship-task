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
