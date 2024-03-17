package http

import (
	"encoding/json"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

const (
	defaultListActorsLimit = 100
	defaultListOffset      = 0
)

type listActorsResponseBody struct {
	Actors []*actor `json:"actors"`
}

// @Summary		List all actors with optional query parameters 'limit' and 'offset'
// @Security		BasicAuth
// @Tags			actors
// @Description	List all actors with optional query parameters 'limit' and 'offset'
// @ID				list-actors
// @Produce		json
// @Param			limit	query		integer	false	"An optional query parameter 'limit' that limits total number of returned actors. By default 'limit' = 100"
// @Param			offset	query		integer	false	"An optional query parameter 'offset' that indicates how many records should be skipped while listing actors. By default 'offset' = 0"
// @Success		200		{object}	listActorsResponseBody
// @Failure		400		{object}	apiv1.Response
// @Failure		401		{object}	apiv1.Response
// @Failure		403		{object}	apiv1.Response
// @Failure		500		{object}	apiv1.Response
// @Router			/actors [get]
func (h *Handler) ListActorsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var limit, offset int
		var err error

		limitStr := request.URL.Query().Get("limit")

		if limitStr == "" {
			limit = defaultListActorsLimit
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			if limit < 0 {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": "'limit' must not be negative"}))

				return
			}
		}

		offsetStr := request.URL.Query().Get("offset")

		if offsetStr == "" {
			offset = defaultListOffset
		} else {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			if offset < 0 {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": "'offset' must not be negative"}))

				return
			}
		}

		domainActors, err := h.actorService.List(limit, offset)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		payload, _ := json.Marshal(&listActorsResponseBody{
			Actors: buildActors(domainActors),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
