package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultListFilmsLimit  = 100
	defaultListFilmsOffset = 0
)

const (
	sortParamTitle       = "title"
	sortParamReleaseDate = "release-date"
	sortParamRating      = "rating"
)

type listFilmsResponseBody struct {
	Films []*film `json:"films"`
}

// @Summary		List all films with optional 'sort', 'limit', 'offset' query parameters
// @Security		BasicAuth
// @Tags			films
// @Description	List all films with the possibility of sorting via 'sort' parameter by 'title' and/or 'rating' and/or 'release-date' and/or with 'limit' and/or 'offset' parameters
// @ID				list-films
// @Produce		json
// @Param			sort	query		string	false	"An optional query parameter 'sort' that indicates how films should be sorted. By default 'sort' = 'rating:desc'. Expected as `title:asc,release-date:desc,rating:desc` in any order of necessary parameters"
// @Param			limit	query		integer	false	"An optional query parameter 'limit' that limits total number of returned films. By default 'limit' = 100"
// @Param			offset	query		integer	false	"An optional query parameter 'offset' that indicates how many records should be skipped while listing films. By default 'offset' = 0"
// @Success		200		{object}	listFilmsResponseBody
// @Failure		400		{object}	apiv1.Response
// @Failure		401		{object}	apiv1.Response
// @Failure		403		{object}	apiv1.Response
// @Failure		500		{object}	apiv1.Response
// @Router			/films [get]
func (h *Handler) ListFilmsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		const handlerName = "ListFilmsHandler"

		log := h.logger.With(slog.String("handlerName", handlerName))

		var limit, offset int
		var err error

		limitStr := request.URL.Query().Get("limit")

		if limitStr == "" {
			limit = defaultListFilmsLimit
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				log.Error("failed to convert 'limit' parameter", "limit", limitStr, "error", err.Error())

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
			offset = defaultListFilmsOffset
		} else {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				log.Error("failed to convert 'offset' parameter", "offset", offsetStr, "error", err.Error())

				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			if offset < 0 {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": "'offset' must not be negative"}))

				return
			}
		}

		var titleOrder, releaseDateOrder, ratingOrder string

		sortQueryParams := request.URL.Query().Get("sort")

		titleOrder, releaseDateOrder, ratingOrder, err = getSortParams(sortQueryParams)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		domainFilms, err := h.filmService.ListWithSort(titleOrder, releaseDateOrder, ratingOrder, limit, offset)
		if err != nil {
			log.Error("failed to list films", "error", err.Error())

			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		payload, _ := json.Marshal(&listFilmsResponseBody{
			Films: buildFilms(domainFilms),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}

func getSortParams(sortParams string) (titleOrder string, releaseDateOrder string, ratingOrder string, err error) {
	if sortParams == "" {
		return "", "", "", nil
	}

	setParamOrder := func(param string, order string) {
		switch param {
		case sortParamTitle:
			titleOrder = order
		case sortParamReleaseDate:
			releaseDateOrder = order
		case sortParamRating:
			ratingOrder = order
		}
	}

	// sortParams expected as 'title:asc,release-date:desc,rating:desc'

	splitSortParams := strings.Split(sortParams, ",")

	for _, param := range splitSortParams {

		// Expect parameter like 'title:asc'
		paramWithSortOrder := strings.Split(param, ":")

		if len(paramWithSortOrder) != 2 {
			return titleOrder, releaseDateOrder, ratingOrder, errors.New("invalid sort parameters. Must be like 'title:asc,release-date:desc,rating:desc'")
		}

		if paramWithSortOrder[0] != sortParamTitle && paramWithSortOrder[0] != sortParamReleaseDate && paramWithSortOrder[0] != sortParamRating {
			return titleOrder, releaseDateOrder, ratingOrder, errors.New(fmt.Sprintf("unexpected sort parameter: '%s'", paramWithSortOrder[0]))
		}

		setParamOrder(paramWithSortOrder[0], paramWithSortOrder[1])
	}

	return titleOrder, releaseDateOrder, ratingOrder, nil
}
