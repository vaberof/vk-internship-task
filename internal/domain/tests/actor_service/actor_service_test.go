package domain_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/vaberof/vk-internship-task/internal/domain"
	mocks "github.com/vaberof/vk-internship-task/internal/domain/mocks"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	actorName := domain.ActorName("Actor")
	actorSex := domain.ActorSex(0)
	actorBirthdate := domain.ActorBirthDate(time.Now())

	expected := &domain.Actor{
		Id:        domain.ActorId(1),
		Name:      actorName,
		Sex:       actorSex,
		BirthDate: actorBirthdate,
		Films:     []*domain.Film{},
	}

	actorStorage.EXPECT().Create(actorName, actorSex, actorBirthdate).Return(expected, nil).Times(1)

	actorService := domain.NewActorService(actorStorage, logsBuilder)
	actor, err := actorService.Create(actorName, actorSex, actorBirthdate)
	require.NoError(t, err)
	require.Equal(t, expected, actor)
}

func TestCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	actorName := domain.ActorName("Actor")
	actorSex := domain.ActorSex(0)
	parsedActorBirthdate, _ := time.Parse(time.DateOnly, "2024-03-18")
	actorBirthdate := domain.ActorBirthDate(parsedActorBirthdate)

	fakeError := errors.New("database is down")
	expectedErr := fmt.Errorf("failed to create an actor: %w", fakeError)

	actorStorage.EXPECT().Create(actorName, actorSex, actorBirthdate).Return(nil, expectedErr).Times(1)

	actorService := domain.NewActorService(actorStorage, logsBuilder)
	actor, err := actorService.Create(actorName, actorSex, actorBirthdate)
	require.Error(t, err)
	require.EqualError(t, expectedErr, err.Error())
	require.Nil(t, actor)
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	actorId := domain.ActorId(1)
	actorName := domain.ActorName("Actor")
	actorSex := domain.ActorSex(0)
	actorBirthdate := domain.ActorBirthDate(time.Now())

	expected := &domain.Actor{
		Id:        actorId,
		Name:      actorName,
		Sex:       actorSex,
		BirthDate: actorBirthdate,
		Films: []*domain.Film{
			{
				Id:          domain.FilmId(1),
				Title:       domain.FilmTitle("Title_1"),
				Description: domain.FilmDescription("Description_1"),
				ReleaseDate: domain.FilmReleaseDate(time.Now()),
				Rating:      domain.FilmRating(10),
			},
			{
				Id:          domain.FilmId(2),
				Title:       domain.FilmTitle("Title_2"),
				Description: domain.FilmDescription("Description_2"),
				ReleaseDate: domain.FilmReleaseDate(time.Now().Add(2 * time.Minute)),
				Rating:      domain.FilmRating(9),
			},
		},
	}

	actorStorage.EXPECT().IsExists(actorId).Return(true, nil).Times(1)
	actorStorage.EXPECT().Update(actorId, &actorName, &actorSex, &actorBirthdate).Return(expected, nil).Times(1)

	actorService := domain.NewActorService(actorStorage, logsBuilder)
	actor, err := actorService.Update(actorId, &actorName, &actorSex, &actorBirthdate)
	require.NoError(t, err)
	require.Equal(t, expected, actor)
}

func TestUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	actorService := domain.NewActorService(actorStorage, logsBuilder)

	type in struct {
		Id        domain.ActorId
		Name      *domain.ActorName
		Sex       *domain.ActorSex
		Birthdate *domain.ActorBirthDate
	}

	var (
		actorId1 = domain.ActorId(1)
		actorId2 = domain.ActorId(2)

		actorName1 = domain.ActorName("Name_1")
		actorName2 = domain.ActorName("Name_2")

		actorSex1 = domain.ActorSex(0)
		actorSex2 = domain.ActorSex(1)

		actorBirthdate1 = domain.ActorBirthDate(time.Now())
		actorBirthdate2 = domain.ActorBirthDate(time.Now().Add(1 * time.Minute))
	)

	testCases := []struct {
		name string
		in   struct {
			Id        domain.ActorId
			Name      *domain.ActorName
			Sex       *domain.ActorSex
			Birthdate *domain.ActorBirthDate
		}
		out          *domain.Actor
		updateExpErr error
		exists       bool
		existsExpErr error
	}{
		{
			name: "err_actor_not_found",
			in: in{
				Id:        actorId1,
				Name:      &actorName1,
				Sex:       &actorSex1,
				Birthdate: &actorBirthdate1,
			},
			out:          nil,
			updateExpErr: domain.ErrActorNotFound,
			exists:       false,
			existsExpErr: nil,
		},
		{
			name: "err_other",
			in: in{
				Id:        actorId2,
				Name:      &actorName2,
				Sex:       &actorSex2,
				Birthdate: &actorBirthdate2,
			},
			out:          nil,
			updateExpErr: fmt.Errorf("failed to update an actor: %w", errors.New("database is down")),
			exists:       true,
			existsExpErr: nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actorStorage.EXPECT().IsExists(tCase.in.Id).Return(tCase.exists, tCase.existsExpErr).AnyTimes()
			actorStorage.EXPECT().Update(tCase.in.Id, tCase.in.Name, tCase.in.Sex, tCase.in.Birthdate).Return(tCase.out, tCase.updateExpErr).AnyTimes()
			actor, err := actorService.Update(tCase.in.Id, tCase.in.Name, tCase.in.Sex, tCase.in.Birthdate)
			require.Error(t, err)
			require.EqualError(t, tCase.updateExpErr, err.Error())
			require.Nil(t, actor)
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	actorId := domain.ActorId(1)

	actorStorage.EXPECT().IsExists(actorId).Return(true, nil).Times(1)
	actorStorage.EXPECT().Delete(actorId).Return(nil).Times(1)

	actorService := domain.NewActorService(actorStorage, logsBuilder)
	err := actorService.Delete(actorId)
	require.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	actorService := domain.NewActorService(actorStorage, logsBuilder)

	testCases := []struct {
		name         string
		in           domain.ActorId
		out          *domain.Actor
		deleteExpErr error
		exists       bool
		existsExpErr error
	}{
		{
			name:         "err_actor_not_found",
			in:           domain.ActorId(1),
			out:          nil,
			deleteExpErr: domain.ErrActorNotFound,
			exists:       false,
			existsExpErr: nil,
		},
		{
			name:         "err_other",
			in:           domain.ActorId(2),
			out:          nil,
			deleteExpErr: fmt.Errorf("failed to delete an actor: %w", errors.New("database is down")),
			exists:       true,
			existsExpErr: nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actorStorage.EXPECT().IsExists(tCase.in).Return(tCase.exists, tCase.existsExpErr).AnyTimes()
			actorStorage.EXPECT().Delete(tCase.in).Return(tCase.deleteExpErr).AnyTimes()
			err := actorService.Delete(tCase.in)
			require.EqualError(t, err, tCase.deleteExpErr.Error())
		})
	}
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	expected := []*domain.Actor{
		{
			Id:        domain.ActorId(1),
			Name:      domain.ActorName("Actor_1"),
			Sex:       domain.ActorSex(0),
			BirthDate: domain.ActorBirthDate(time.Now()),
			Films: []*domain.Film{
				{
					Id:          domain.FilmId(1),
					Title:       domain.FilmTitle("Title_1"),
					Description: domain.FilmDescription("Description_1"),
					ReleaseDate: domain.FilmReleaseDate(time.Now()),
					Rating:      domain.FilmRating(10),
				},
				{
					Id:          domain.FilmId(2),
					Title:       domain.FilmTitle("Title_2"),
					Description: domain.FilmDescription("Description_2"),
					ReleaseDate: domain.FilmReleaseDate(time.Now().Add(2 * time.Minute)),
					Rating:      domain.FilmRating(9),
				},
				{
					Id:          domain.FilmId(3),
					Title:       domain.FilmTitle("Title_3"),
					Description: domain.FilmDescription("Description_3"),
					ReleaseDate: domain.FilmReleaseDate(time.Now().Add(3 * time.Minute)),
					Rating:      domain.FilmRating(8),
				},
			},
		},
		{
			Id:        domain.ActorId(2),
			Name:      domain.ActorName("Actor_2"),
			Sex:       domain.ActorSex(1),
			BirthDate: domain.ActorBirthDate(time.Now().Add(1 * time.Minute)),
			Films: []*domain.Film{
				{
					Id:          domain.FilmId(1),
					Title:       domain.FilmTitle("Title_1"),
					Description: domain.FilmDescription("Description_1"),
					ReleaseDate: domain.FilmReleaseDate(time.Now()),
					Rating:      domain.FilmRating(10),
				},
				{
					Id:          domain.FilmId(2),
					Title:       domain.FilmTitle("Title_2"),
					Description: domain.FilmDescription("Description_2"),
					ReleaseDate: domain.FilmReleaseDate(time.Now().Add(2 * time.Minute)),
					Rating:      domain.FilmRating(9),
				},
			},
		},
	}

	limit := 100
	offset := 0

	actorStorage.EXPECT().List(limit, offset).Return(expected, nil).Times(1)

	actorService := domain.NewActorService(actorStorage, logsBuilder)
	actors, err := actorService.List(limit, offset)
	require.NoError(t, err)
	require.Equal(t, expected, actors)
}

func TestListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actorStorage := mocks.NewMockActorStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	fakeError := errors.New("database is down")
	expectedErr := fmt.Errorf("failed to list actors: %w", fakeError)

	limit := 100
	offset := 0

	actorStorage.EXPECT().List(limit, offset).Return(nil, expectedErr).Times(1)

	actorService := domain.NewActorService(actorStorage, logsBuilder)
	actors, err := actorService.List(limit, offset)
	require.Error(t, err)
	require.EqualError(t, expectedErr, err.Error())
	require.Nil(t, actors)
}
