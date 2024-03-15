package postgres

import "time"

type PgActor struct {
	Id        int64
	Name      string
	Sex       uint8
	BirthDate time.Time
	Films     []*PgFilm
}
