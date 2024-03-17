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

type deleteActorResponseBody struct {
	Message string `json:"message"`
}

func (h *Handler) DeleteActorHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		actorIdPathParam := request.PathValue("id")
		actorId, err := strconv.Atoi(actorIdPathParam)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		err = h.actorService.Delete(domain.ActorId(actorId))
		if err != nil {
			if errors.Is(err, domain.ErrActorNotFound) {
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageActorNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))
			}

			return
		}

		payload, _ := json.Marshal(&deleteActorResponseBody{
			Message: fmt.Sprintf("Actor with id '%d' has deleted successfully", actorId),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
