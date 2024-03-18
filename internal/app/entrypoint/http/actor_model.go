package http

import (
	"github.com/vaberof/vk-internship-task/internal/domain"
	"time"
)

type actor struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	Sex       uint8        `json:"sex"`
	BirthDate string       `json:"birthdate"`
	Films     []*actorFilm `json:"films,omitempty"`
}

type actorFilm struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Rating      uint8  `json:"rating"`
}

func buildDomainActorIds(actorIds []int64) []domain.ActorId {
	domainActorIds := make([]domain.ActorId, len(actorIds))
	for i := range actorIds {
		domainActorIds[i] = domain.ActorId(actorIds[i])
	}
	return domainActorIds
}

func buildActors(domainActors []*domain.Actor) []*actor {
	actors := make([]*actor, len(domainActors))
	for i := range domainActors {
		actors[i] = buildActor(domainActors[i])
	}
	return actors
}

func buildActor(domainActor *domain.Actor) *actor {
	return &actor{
		Id:        domainActor.Id.Int64(),
		Name:      domainActor.Name.String(),
		Sex:       domainActor.Sex.Uint8(),
		BirthDate: domainActor.BirthDate.Time().Format(time.DateOnly),
		Films:     buildActorFilms(domainActor.Films),
	}
}

func buildActorFilms(domainFilms []*domain.Film) []*actorFilm {
	actorFilms := make([]*actorFilm, len(domainFilms))
	for i := range domainFilms {
		actorFilms[i] = buildActorFilm(domainFilms[i])
	}
	return actorFilms
}

func buildActorFilm(domainFilm *domain.Film) *actorFilm {
	return &actorFilm{
		Id:          domainFilm.Id.Int64(),
		Title:       domainFilm.Title.String(),
		Description: domainFilm.Description.String(),
		ReleaseDate: domainFilm.ReleaseDate.Time().Format(time.DateOnly),
		Rating:      domainFilm.Rating.Uint8(),
	}
}
