package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"log/slog"
	"net/http"
	"strconv"
)

type deleteFilmResponseBody struct {
	Message string `json:"message"`
}

// @Summary		Delete a film by path parameter 'id'
// @Security		BasicAuth
// @Tags			films
// @Description	Delete a film by path parameter 'id'
// @ID				delete-film
// @Accept			json
// @Produce		json
// @Param			id	path		integer	true	"Film`s id that needs to be deleted"
// @Success		200	{object}	deleteFilmResponseBody
// @Failure		400	{object}	apiv1.Response
// @Failure		401	{object}	apiv1.Response
// @Failure		403	{object}	apiv1.Response
// @Failure		404	{object}	apiv1.Response
// @Failure		500	{object}	apiv1.Response
// @Router			/films/{id} [delete]
func (h *Handler) DeleteFilmHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		const handlerName = "DeleteFilmHandler"

		log := h.logger.With(slog.String("handlerName", handlerName))

		filmIdPathParam := request.PathValue("id")
		if filmIdPathParam == "" {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageFilmInvalidRequestBody, apiv1.ErrorDescription{"error": "missing required path parameter 'id'"}))

			return
		}

		filmId, err := strconv.Atoi(filmIdPathParam)
		if err != nil {
			log.Error("failed to convert 'id' parameter", "id", filmIdPathParam, "error", err.Error())

			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		err = h.filmService.Delete(domain.FilmId(filmId))
		if err != nil {
			if errors.Is(err, domain.ErrFilmNotFound) {
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageFilmNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
				log.Error("failed to delete film", "id", filmId, "error", err.Error())

				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))
			}

			return
		}

		payload, _ := json.Marshal(&deleteFilmResponseBody{
			Message: fmt.Sprintf("Film with id '%d' has deleted successfully", filmId),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
