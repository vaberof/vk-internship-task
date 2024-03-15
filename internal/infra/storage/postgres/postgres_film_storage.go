package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/vk-internship-task/internal/domain"
	"github.com/vaberof/vk-internship-task/internal/infra/storage"
	"strings"
)

type PgFilmStorage struct {
	db *sqlx.DB
}

func NewPgFilmStorage(db *sqlx.DB) *PgFilmStorage {
	return &PgFilmStorage{db: db}
}

func (s *PgFilmStorage) Create(title domain.FilmTitle, description domain.FilmDescription, releaseDate domain.FilmReleaseDate, rating domain.FilmRating, actorIds []domain.ActorId) (*domain.Film, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while creating film: %w", err)
	}
	defer tx.Rollback()

	var film PgFilm

	query := `
			INSERT INTO films (
			                    title,
			                    description,
			                    release_date,
			                   	rating
				) VALUES ($1, $2, $3, $4)
				RETURNING 
					id, 
				 	title,
					description,
					release_date,
					rating
`

	row := tx.QueryRow(query, title, description, releaseDate, rating)
	if err = row.Scan(
		&film.Id,
		&film.Title,
		&film.Description,
		&film.ReleaseDate,
		&film.Rating,
	); err != nil {
		return nil, fmt.Errorf("failed to create a film: %w", err)
	}

	queryFilmsActors := `
						INSERT INTO films_actors(film_id, actor_id)
						VALUES ($1, $2)
`

	for _, actorId := range actorIds {
		_, err = tx.Exec(queryFilmsActors, film.Id, actorId)
		if err != nil {
			return nil, fmt.Errorf("failed to create a film: failed to insert values to 'films_actors' table %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while creating film: %w", err)
	}

	return buildDomainFilm(&film), nil
}

func (s *PgFilmStorage) Update(id domain.FilmId, title *domain.FilmTitle, description *domain.FilmDescription, releaseDate *domain.FilmReleaseDate, rating *domain.FilmRating) (*domain.Film, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while updating film: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("LOCK TABLE films IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		return nil, fmt.Errorf("failed to lock 'films' table while updating film: %w", err)
	}

	var film PgFilm

	query := `
			UPDATE films 
						SET title=COALESCE($1, title),
							description=COALESCE($2, description),
							release_date=COALESCE($3, release_date),
							rating=COALSECE($4, rating)
					 	WHERE id=$5
			RETURNING
			    id, 
				title,
				description,
				release_date,
				rating
`

	row := tx.QueryRow(query, title, description, releaseDate, rating, id)
	if err = row.Scan(
		&film.Id,
		&film.Title,
		&film.Description,
		&film.ReleaseDate,
		&film.Rating,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update a film: %w", storage.ErrFilmNotFound)
		}
		return nil, fmt.Errorf("failed to update a film: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while updating film: %w", err)
	}

	return buildDomainFilm(&film), nil
}

func (s *PgFilmStorage) Delete(id domain.FilmId) error {
	query := `DELETE FROM films WHERE id=$1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete film: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("failed to delete film: %w", storage.ErrFilmNotFound)
	}
	return nil
}

func (s *PgFilmStorage) ListWithSort(title *domain.FilmTitle, releaseDate *domain.FilmReleaseDate, rating *domain.FilmRating, limit, offset int) ([]*domain.Film, error) {
	limitOffsetParams := fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)
	orderParam := s.buildOrderParam(title, releaseDate, rating)

	query := `
			SELECT f.id,
			       f.title,
			       f.description,
			       f.release_date,
			       f.rating,
			       a.id,
			       a.name,
			       a.sex,
			       a.birthdate
			FROM films AS f
			JOIN films_actors AS fa ON f.id = fa.film_id
			JOIN actors AS a ON a.id = fa.actor_id
` + orderParam + limitOffsetParams

	var films []*PgFilm

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list films: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var film PgFilm
		var actor PgActor

		if err = rows.Scan(
			&film.Id,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
			&actor.Id,
			&actor.Name,
			&actor.Sex,
			&actor.BirthDate,
		); err != nil {
			return nil, fmt.Errorf("failed to list films: %w", err)
		}

		var filmExists bool
		for _, pgFilm := range films {
			if pgFilm.Id == film.Id {
				pgFilm.Actors = append(pgFilm.Actors, &actor)
				filmExists = true
				break
			}
		}
		if !filmExists {
			film.Actors = append(film.Actors, &actor)
			films = append(films, &film)
		}
	}

	return buildDomainFilms(films), nil
}

func (s *PgFilmStorage) SearchByFilters(title domain.FilmTitle, actorName domain.ActorName) ([]*domain.Film, error) {
	query := `
			SELECT f.id,
			       f.title,
			       f.description,
			       f.release_date,
			       f.rating,
			       a.id,
			       a.name,
			       a.sex,
			       a.birthdate
			FROM films AS f
			JOIN films_actors AS fa ON f.id = fa.film_id
			JOIN actors AS a ON a.id = fa.actor_id
			WHERE f.title LIKE '%' || $1 || '%' AND a.name LIKE '%' || $2 || '%'
`

	var films []*PgFilm

	rows, err := s.db.Query(query, title, actorName)
	if err != nil {
		return nil, fmt.Errorf("failed to list films: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var film PgFilm
		var actor PgActor

		if err = rows.Scan(
			&film.Id,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
			&actor.Id,
			&actor.Name,
			&actor.Sex,
			&actor.BirthDate,
		); err != nil {
			return nil, fmt.Errorf("failed to list films: %w", err)
		}

		var filmExists bool
		for _, pgFilm := range films {
			if pgFilm.Id == film.Id {
				pgFilm.Actors = append(pgFilm.Actors, &actor)
				filmExists = true
				break
			}
		}
		if !filmExists {
			film.Actors = append(film.Actors, &actor)
			films = append(films, &film)
		}
	}

	return buildDomainFilms(films), nil
}

func (s *PgFilmStorage) IsExists(id domain.FilmId) (bool, error) {
	query := `
			SELECT id FROM films 
			WHERE id=$1
`
	var filmdId int64
	err := s.db.QueryRow(query, id).Scan(&filmdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check whether film exists or not: %w", err)
	}
	return true, nil
}

func (s *PgFilmStorage) buildOrderParam(title *domain.FilmTitle, releaseDate *domain.FilmReleaseDate, rating *domain.FilmRating) string {
	orderParam := " ORDER BY "

	var orderByClauses []string

	if title != nil {
		orderByClauses = append(orderByClauses, "f.title")
	}
	if releaseDate != nil {
		orderByClauses = append(orderByClauses, "f.release_date")
	}
	if rating != nil {
		orderByClauses = append(orderByClauses, "f.rating DESC")
	}

	if len(orderByClauses) > 0 {
		orderParam += strings.Join(orderByClauses, ",") + " "
	} else {
		orderParam += "f.rating DESC "
	}

	return orderParam
}

func buildDomainFilms(postgresFilms []*PgFilm) []*domain.Film {
	domainFilms := make([]*domain.Film, 0, len(postgresFilms))
	for i := range postgresFilms {
		domainFilms[i] = buildDomainFilm(postgresFilms[i])
	}
	return domainFilms
}

func buildDomainFilm(postgresFilm *PgFilm) *domain.Film {
	return &domain.Film{
		Id:          domain.FilmId(postgresFilm.Id),
		Title:       domain.FilmTitle(postgresFilm.Title),
		Description: domain.FilmDescription(postgresFilm.Description.String),
		ReleaseDate: domain.FilmReleaseDate(postgresFilm.ReleaseDate),
		Rating:      domain.FilmRating(postgresFilm.Rating),
		Actors:      buildDomainActors(postgresFilm.Actors),
	}
}
