package domain

import (
	"errors"
	"fmt"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log/slog"
)

var (
	ErrActorNotFound  = errors.New("actor not found")
	ErrActorsNotFound = errors.New("actors not found")
)

type ActorService interface {
	Create(name ActorName, sex ActorSex, birthDate ActorBirthDate) (*Actor, error)
	Update(id ActorId, name *ActorName, sex *ActorSex, birthDate *ActorBirthDate) (*Actor, error)
	Delete(id ActorId) error
	List(limit, offset int) ([]*Actor, error)
}

type actorServiceImpl struct {
	actorStorage ActorStorage

	logger *slog.Logger
}

func NewActorService(actorStorage ActorStorage, logsBuilder *logs.Logs) ActorService {
	logger := logsBuilder.WithName("domain.service.actor")
	return &actorServiceImpl{
		actorStorage: actorStorage,
		logger:       logger,
	}
}

func (a *actorServiceImpl) Create(name ActorName, sex ActorSex, birthDate ActorBirthDate) (*Actor, error) {
	const operation = "Create"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("name", name.String()))

	log.Info("creating an actor")

	domainActor, err := a.actorStorage.Create(name, sex, birthDate)
	if err != nil {
		log.Error("failed to create an actor", "error", err)
		return nil, err
	}

	log.Info("actor has created")

	return domainActor, nil
}

func (a *actorServiceImpl) Update(id ActorId, name *ActorName, sex *ActorSex, birthDate *ActorBirthDate) (*Actor, error) {
	const operation = "Update"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.Int64("id", id.Int64()))

	log.Info("updating an actor")

	exists, err := a.actorStorage.IsExists(id)
	if err != nil {
		log.Error("failed to update an actor", "error", err)
		return nil, err
	}
	if !exists {
		log.Warn("failed to update an actor", "error", fmt.Sprintf("actor with id '%d' not found", id.Int64()))
		return nil, ErrActorNotFound
	}

	domainActor, err := a.actorStorage.Update(id, name, sex, birthDate)
	if err != nil {
		log.Error("failed to update an actor", "error", err)
		return nil, err
	}

	log.Info("actor has updated")

	return domainActor, nil
}

func (a *actorServiceImpl) Delete(id ActorId) error {
	const operation = "Delete"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.Int64("id", id.Int64()))

	log.Info("deleting an actor")

	exists, err := a.actorStorage.IsExists(id)
	if err != nil {
		log.Error("failed to delete an actor", "error", err)
		return err
	}
	if !exists {
		log.Warn("failed to delete an actor", "error", fmt.Sprintf("actor with id '%d' not found", id.Int64()))
		return ErrActorNotFound
	}

	err = a.actorStorage.Delete(id)
	if err != nil {
		log.Error("failed to delete an actor", "error", err)
		return err
	}

	log.Info("actor has deleted")

	return nil
}

func (a *actorServiceImpl) List(limit, offset int) ([]*Actor, error) {
	const operation = "List"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.Int("limit", limit),
		slog.Int("offset", offset),
	)

	log.Info("listing actors")

	domainActors, err := a.actorStorage.List(limit, offset)
	if err != nil {
		log.Error("failed to list actors", "error", err)
		return nil, err
	}

	log.Info("actors have listed")

	return domainActors, nil
}
