package postgres

import (
	"database/sql"
	"time"
)

type PgFilm struct {
	Id          int64
	Title       string
	Description sql.NullString
	ReleaseDate time.Time
	Rating      uint8
	Actors      []*PgActor
}
