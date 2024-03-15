package storage

import "errors"

var (
	ErrActorNotFound = errors.New("actor not found")
	ErrFilmNotFound  = errors.New("film not found")
)
