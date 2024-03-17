package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrActorNotFound = errors.New("actor not found")
	ErrFilmNotFound  = errors.New("film not found")
)
