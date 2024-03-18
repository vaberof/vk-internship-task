package http

import (
	"github.com/vaberof/vk-internship-task/internal/domain"
	"time"
)

type film struct {
	Id          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	ReleaseDate string       `json:"release_date"`
	Rating      uint8        `json:"rating"`
	Actors      []*filmActor `json:"actors,omitempty"`
}

type filmActor struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Sex       uint8  `json:"sex"`
	BirthDate string `json:"birthdate"`
}

func buildFilms(domainFilms []*domain.Film) []*film {
	films := make([]*film, len(domainFilms))
	for i := range domainFilms {
		films[i] = buildFilm(domainFilms[i])
	}
	return films
}

func buildFilm(domainFilm *domain.Film) *film {
	return &film{
		Id:          domainFilm.Id.Int64(),
		Title:       domainFilm.Title.String(),
		Description: domainFilm.Description.String(),
		ReleaseDate: domainFilm.ReleaseDate.Time().Format(time.DateOnly),
		Rating:      domainFilm.Rating.Uint8(),
		Actors:      buildFilmActors(domainFilm.Actors),
	}
}

func buildFilmActors(domainActors []*domain.Actor) []*filmActor {
	filmActors := make([]*filmActor, len(domainActors))
	for i := range domainActors {
		filmActors[i] = buildFilmActor(domainActors[i])
	}
	return filmActors
}

func buildFilmActor(domainActor *domain.Actor) *filmActor {
	return &filmActor{
		Id:        domainActor.Id.Int64(),
		Name:      domainActor.Name.String(),
		Sex:       domainActor.Sex.Uint8(),
		BirthDate: domainActor.BirthDate.Time().Format(time.DateOnly),
	}
}
