package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

type deleteFilmResponseBody struct {
	Message string `json:"message"`
}

func (h *Handler) DeleteFilmHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		filmIdPathParam := request.PathValue("id")
		filmId, err := strconv.Atoi(filmIdPathParam)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageFilmInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		err = h.filmService.Delete(domain.FilmId(filmId))
		if err != nil {
			if errors.Is(err, domain.ErrFilmNotFound) {
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageFilmNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
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
