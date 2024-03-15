package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/internal/infra/storage"
	"strconv"
	"strings"
)

type PgActorStorage struct {
	db *sqlx.DB
}

func NewPgActorStorage(db *sqlx.DB) *PgActorStorage {
	return &PgActorStorage{db: db}
}

func (s *PgActorStorage) Create(name domain.ActorName, sex domain.ActorSex, birthDate domain.ActorBirthDate) (*domain.Actor, error) {
	var actor PgActor
	query := `
			INSERT INTO actors (
			                    name,
			                    sex,
			                    birthdate
				) VALUES ($1, $2, $3)
				RETURNING 
					id, 
					name,
					sex, birthdate
`
	row := s.db.QueryRow(query, name, sex, birthDate)
	if err := row.Scan(
		&actor.Id,
		&actor.Name,
		&actor.Sex,
		&actor.BirthDate,
	); err != nil {
		return nil, fmt.Errorf("failed to create an actor: %w", err)
	}
	return buildDomainActor(&actor), nil
}

func (s *PgActorStorage) Update(id domain.ActorId, name *domain.ActorName, sex *domain.ActorSex, birthDate *domain.ActorBirthDate) (*domain.Actor, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while updating actor: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("LOCK TABLE actors IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		return nil, fmt.Errorf("failed to lock 'actors' table while updating actor: %w", err)
	}

	var actor PgActor

	query := `
		UPDATE actors
					SET	name=COALESCE($1, name),
						sex=COALESCE($2, sex),
						birthdate=COALESCE($3, birthdate)
					WHERE id=$4
		RETURNING
			id,
			name,
			sex,
			birthdate
`

	row := tx.QueryRow(query, name, sex, birthDate, id)

	if err = row.Scan(
		&actor.Id,
		&actor.Name,
		&actor.Sex,
		&actor.BirthDate,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update actor in database: %w", storage.ErrActorNotFound)
		}
		return nil, fmt.Errorf("failed to update actor in database: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while updating actor: %w", err)
	}

	return buildDomainActor(&actor), nil
}

func (s *PgActorStorage) Delete(id domain.ActorId) error {
	query := `DELETE FROM actors WHERE id=$1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("failed to delete actor: %w", storage.ErrActorNotFound)
	}
	return nil
}

func (s *PgActorStorage) List(limit, offset int) ([]*domain.Actor, error) {
	limitOffsetParams := fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	query := `
			SELECT a.id,
			       a.name,
			       a.sex,
			       a.birthdate,
			       f.id,
			       f.title,
			       f.description,
			       f.release_date,
			       f.rating
		    FROM actors AS a
		    INNER JOIN films_actors AS fa ON a.id = fa.actor_id
		    INNER JOIN films AS f ON f.id = fa.film_id
			` + limitOffsetParams

	var actors []*PgActor

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list actors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var actor PgActor
		var film PgFilm

		if err = rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.Sex,
			&actor.BirthDate,
			&film.Id,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		); err != nil {
			return nil, fmt.Errorf("failed to list actors: %w", err)
		}

		var actorExists bool
		for _, pgActor := range actors {
			if pgActor.Id == actor.Id {
				pgActor.Films = append(pgActor.Films, &film)
				actorExists = true
				break
			}
		}
		if !actorExists {
			actor.Films = append(actor.Films, &film)
			actors = append(actors, &actor)
		}
	}

	return buildDomainActors(actors), nil
}

func (s *PgActorStorage) IsExists(id domain.ActorId) (bool, error) {
	query := `
			SELECT id FROM actors 
			WHERE id=$1
`
	var actorId int64
	err := s.db.QueryRow(query, id).Scan(&actorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check whether actor exists or not: %w", err)
	}
	return true, nil
}

func (s *PgActorStorage) AreExists(ids []domain.ActorId) (bool, error) {
	query := `
			SELECT COUNT(*) FROM actors 
			WHERE id IN ($1)
`
	strIds := make([]string, len(ids))
	for i := range ids {
		strIds[i] = strconv.Itoa(int(ids[i].Int64()))
	}

	strIdsWithComma := strings.Join(strIds, ",")

	var idsCount int64

	err := s.db.QueryRow(query, strIdsWithComma).Scan(&idsCount)
	if err != nil {
		return false, fmt.Errorf("failed to check whether actors exists or not: %w", err)
	}

	if int(idsCount) != len(ids) {
		return false, nil
	}

	return true, nil
}

func buildDomainActors(postgresActors []*PgActor) []*domain.Actor {
	domainActors := make([]*domain.Actor, len(postgresActors))
	for i := range postgresActors {
		domainActors[i] = buildDomainActor(postgresActors[i])
	}
	return domainActors
}

func buildDomainActor(postgresActor *PgActor) *domain.Actor {
	return &domain.Actor{
		Id:        domain.ActorId(postgresActor.Id),
		Name:      domain.ActorName(postgresActor.Name),
		Sex:       domain.ActorSex(postgresActor.Sex),
		BirthDate: domain.ActorBirthDate(postgresActor.BirthDate),
	}
}
