package domain

import (
	"time"
)

type FilmId int64

func (filmId *FilmId) Int64() int64 {
	return int64(*filmId)
}

type FilmTitle string

func (filmTitle *FilmTitle) String() string {
	return string(*filmTitle)
}

type FilmDescription string

func (filmDescription *FilmDescription) String() string {
	return string(*filmDescription)
}

type FilmReleaseDate time.Time

func (filmReleaseDate *FilmReleaseDate) Time() time.Time {
	return time.Time(*filmReleaseDate)
}

type FilmRating uint8

func (filmRating *FilmRating) Uint8() uint8 {
	return uint8(*filmRating)
}

type Film struct {
	Id          FilmId
	Title       FilmTitle
	Description FilmDescription
	ReleaseDate FilmReleaseDate
	Rating      FilmRating
	Actors      []*Actor
}
