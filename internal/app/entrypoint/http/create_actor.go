package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http/views"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
	"time"
)

type createActorRequestBody struct {
	Name      string `json:"name" validate:"required,max=100"`
	Sex       uint8  `json:"sex" validate:"required,oneof=0 1 2 9"`
	BirthDate string `json:"birthdate" validate:"required"`
}

type createActorResponseBody struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Sex       uint8   `json:"sex"`
	BirthDate string  `json:"birthdate"`
	Films     []*film `json:"films"`
}

func (h *Handler) CreateActorHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var createActorReqBody createActorRequestBody
		err := json.NewDecoder(request.Body).Decode(&createActorReqBody)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": "invalid request body"}))

			return
		}

		err = h.validator.Struct(&createActorReqBody)
		if err != nil {
			errors, ok := err.(validator.ValidationErrors)
			if !ok {
				views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": "failed to get validation errors"}))
			} else {
				views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": errors.Error()}))
			}

			return
		}

		birthdate, err := time.Parse(time.DateOnly, createActorReqBody.BirthDate)
		if err != nil {
			views.RenderJSON(rw, http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, ErrMessageActorInvalidRequestBody, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		domainActor, err := h.actorService.Create(
			domain.ActorName(createActorReqBody.Name),
			domain.ActorSex(createActorReqBody.Sex),
			domain.ActorBirthDate(birthdate),
		)
		if err != nil {
			views.RenderJSON(rw, http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, ErrMessageActorInternalServerError, apiv1.ErrorDescription{"error": err.Error()}))

			return
		}

		payload, _ := json.Marshal(&createActorResponseBody{
			Id:        domainActor.Id.Int64(),
			Name:      domainActor.Name.String(),
			Sex:       domainActor.Sex.Uint8(),
			BirthDate: domainActor.BirthDate.Time().Format(time.DateOnly),
			Films:     buildFilms(domainActor.Films),
		})

		views.RenderJSON(rw, http.StatusOK, apiv1.Success(payload))
	}
}
