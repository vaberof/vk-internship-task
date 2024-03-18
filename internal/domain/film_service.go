package domain

import (
	"errors"
	"fmt"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log/slog"
)

var (
	ErrFilmNotFound       = errors.New("film not found")
	ErrFilmActorsNotFound = errors.New("actors not found")
)

type FilmService interface {
	Create(title FilmTitle, description FilmDescription, releaseDate FilmReleaseDate, rating FilmRating, actorIds []ActorId) (*Film, error)
	Update(id FilmId, title *FilmTitle, description *FilmDescription, releaseDate *FilmReleaseDate, rating *FilmRating, actorIds *[]ActorId) (*Film, error)
	Delete(id FilmId) error
	ListWithSort(titleOrder, releaseDateOrder, ratingOrder string, limit, offset int) ([]*Film, error)
	SearchByFilters(title FilmTitle, actorName ActorName, limit, offset int) ([]*Film, error)
}

type filmServiceImpl struct {
	filmStorage  FilmStorage
	actorStorage ActorStorage

	logger *slog.Logger
}

func NewFilmService(filmStorage FilmStorage, actorStorage ActorStorage, logsBuilder *logs.Logs) FilmService {
	logger := logsBuilder.WithName("domain.service.film")
	return &filmServiceImpl{
		filmStorage:  filmStorage,
		actorStorage: actorStorage,
		logger:       logger,
	}
}

func (f *filmServiceImpl) Create(title FilmTitle, description FilmDescription, releaseDate FilmReleaseDate, rating FilmRating, actorIds []ActorId) (*Film, error) {
	const operation = "Create"

	log := f.logger.With(
		slog.String("operation", operation),
		slog.String("title", title.String()),
		slog.Any("actorIds", actorIds),
	)

	log.Info("creating a film")

	exists, err := f.actorStorage.AreExists(actorIds)
	if err != nil {
		log.Error("failed to create a film", "error", err)
		return nil, err
	}
	if !exists {
		log.Warn("failed to create a film", "error", fmt.Sprintf("actors with ids '%v' not found", actorIds))
		return nil, ErrFilmActorsNotFound
	}

	domainFilm, err := f.filmStorage.Create(title, description, releaseDate, rating, actorIds)
	if err != nil {
		log.Error("failed to create a film", "error", err)
		return nil, err
	}

	log.Info("film has created")

	return domainFilm, nil
}

func (f *filmServiceImpl) Update(id FilmId, title *FilmTitle, description *FilmDescription, releaseDate *FilmReleaseDate, rating *FilmRating, actorIds *[]ActorId) (*Film, error) {
	const operation = "Update"

	log := f.logger.With(
		slog.String("operation", operation),
		slog.Int64("id", id.Int64()),
	)

	log.Info("updating a film")

	exists, err := f.filmStorage.IsExists(id)
	if err != nil {
		log.Error("failed to update a film", "error", err)
		return nil, err
	}
	if !exists {
		log.Warn("failed to update a film", "error", fmt.Sprintf("film with id '%d' not found", id.Int64()))
		return nil, ErrFilmNotFound
	}

	if actorIds != nil {
		exists, err = f.actorStorage.AreExists(*actorIds)
		if err != nil {
			log.Error("failed to update a film", "error", err)
			return nil, err
		}
		if !exists {
			log.Warn("failed to update a film", "error", fmt.Sprintf("actors with ids '%v' not found", actorIds))
			return nil, ErrFilmActorsNotFound
		}
	}

	domainFilm, err := f.filmStorage.Update(id, title, description, releaseDate, rating, actorIds)
	if err != nil {
		log.Error("failed to update a film", "error", err)
		return nil, err
	}

	log.Info("film has updated")

	return domainFilm, nil
}

func (f *filmServiceImpl) Delete(id FilmId) error {
	const operation = "Delete"

	log := f.logger.With(
		slog.String("operation", operation),
		slog.Int64("id", id.Int64()),
	)

	log.Info("deleting a film")

	exists, err := f.filmStorage.IsExists(id)
	if err != nil {
		log.Error("failed to delete a film", "error", err)
		return err
	}
	if !exists {
		log.Warn("failed to delete a film", "error", fmt.Sprintf("film with id '%d' not found", id.Int64()))
		return ErrFilmNotFound
	}

	err = f.filmStorage.Delete(id)
	if err != nil {
		log.Error("failed to delete a film", "error", err)
		return err
	}

	log.Info("film has deleted")

	return nil
}

func (f *filmServiceImpl) ListWithSort(titleOrder, releaseDateOrder, ratingOrder string, limit, offset int) ([]*Film, error) {
	const operation = "ListWithSort"

	log := f.logger.With(
		slog.String("operation", operation),
		slog.Int("limit", limit),
		slog.Int("offset", offset),
	)

	log.Info("listing films")

	domainFilms, err := f.filmStorage.ListWithSort(titleOrder, releaseDateOrder, ratingOrder, limit, offset)
	if err != nil {
		log.Error("failed to list films", "error", err)
		return nil, err
	}

	log.Info("films have listed")

	return domainFilms, nil
}

func (f *filmServiceImpl) SearchByFilters(title FilmTitle, actorName ActorName, limit, offset int) ([]*Film, error) {
	const operation = "SearchByFilters"

	log := f.logger.With(
		slog.String("operation", operation),
		slog.String("title", title.String()),
		slog.String("actorName", actorName.String()),
	)

	log.Info("searching films")

	domainFilms, err := f.filmStorage.SearchByFilters(title, actorName, limit, offset)
	if err != nil {
		log.Error("failed to search films", "error", err)
		return nil, err
	}

	log.Info("films have searched")

	return domainFilms, nil
}
