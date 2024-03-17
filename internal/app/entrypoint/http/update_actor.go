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

type updateActorRequestBody struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,max=100"`
	Sex       *uint8  `json:"sex,omitempty" validate:"omitempty,oneof=0 1 2 9"`
	BirthDate *string `json:"birthdate,omitempty"`
}

type updateActorResponseBody struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Sex       uint8   `json:"sex"`
	BirthDate string  `json:"birthdate"`
	Films     []*film `json:"films"`
}

func (h *Handler) UpdateActorHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var updateActorReqBody updateActorRequestBody
		err := json.NewDecoder(request.Body).Decode(&updateActorReqBody)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": "invalid request body"}))

			return
		}

		err = h.validator.Struct(&updateActorReqBody)
		if err != nil {
			errors, ok := err.(validator.ValidationErrors)
			if !ok {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": "failed to get validation errors"}))
			} else {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": errors.Error()}))
			}

			return
		}

		var domainBirthdate *domain.ActorBirthDate
		if updateActorReqBody.BirthDate != nil {
			parsedBirthdate, err := time.Parse(time.DateOnly, *updateActorReqBody.BirthDate)
			if err != nil {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": err.Error()}))

				return
			}
			convDomainBirthdate := domain.ActorBirthDate(parsedBirthdate)
			domainBirthdate = &convDomainBirthdate
		}

		actorIdPathParam := request.PathValue("id")
		actorId, err := strconv.Atoi(actorIdPathParam)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		var domainActorName *domain.ActorName
		if updateActorReqBody.Name != nil {
			convDomainActorName := domain.ActorName(*updateActorReqBody.Name)
			domainActorName = &convDomainActorName
		}

		var domainActorSex *domain.ActorSex
		if updateActorReqBody.Sex != nil {
			convDomainActorSex := domain.ActorSex(*updateActorReqBody.Sex)
			domainActorSex = &convDomainActorSex
		}

		domainActor, err := h.actorService.Update(
			domain.ActorId(actorId),
			domainActorName,
			domainActorSex,
			domainBirthdate,
		)
		if err != nil {
			if errors.Is(err, domain.ErrActorNotFound) {
				views.RenderJSON(rw, http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, ErrMessageActorNotFound, apiv1.ErrorDescription{"error": err.Error()}))
			} else {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))
			}

			return
		}

		payload, _ := json.Marshal(&updateActorResponseBody{
			Id:        domainActor.Id.Int64(),
			Name:      domainActor.Name.String(),
			Sex:       domainActor.Sex.Uint8(),
			BirthDate: domainActor.BirthDate.Time().Format(time.DateOnly),
			Films:     buildFilms(domainActor.Films),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
