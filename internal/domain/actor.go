package domain

import (
	"time"
)

type ActorId int64

func (actorId *ActorId) Int64() int64 {
	return int64(*actorId)
}

type ActorName string

func (actorName *ActorName) String() string {
	return string(*actorName)
}

type ActorSex uint8

func (actorSex *ActorSex) Uint8() uint8 {
	return uint8(*actorSex)
}

type ActorBirthDate time.Time

func (actorBirthDate *ActorBirthDate) Time() time.Time {
	return time.Time(*actorBirthDate)
}

type Actor struct {
	Id        ActorId
	Name      ActorName
	Sex       ActorSex
	BirthDate ActorBirthDate
	Films     []*Film
}
