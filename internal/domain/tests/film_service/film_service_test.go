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

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	filmId := domain.FilmId(1)
	filmTitle := domain.FilmTitle("Title_1")
	filmDescription := domain.FilmDescription("Description_1")
	filmReleaseDate := domain.FilmReleaseDate(time.Now())
	filmRating := domain.FilmRating(10)
	actorIds := []domain.ActorId{1, 2}

	expected := &domain.Film{
		Id:          filmId,
		Title:       filmTitle,
		Description: filmDescription,
		ReleaseDate: filmReleaseDate,
		Rating:      filmRating,
		Actors: []*domain.Actor{
			{
				Id:        domain.ActorId(1),
				Name:      domain.ActorName("Actor_1"),
				Sex:       domain.ActorSex(0),
				BirthDate: domain.ActorBirthDate(time.Now()),
			},
			{
				Id:        domain.ActorId(2),
				Name:      domain.ActorName("Actor_2"),
				Sex:       domain.ActorSex(1),
				BirthDate: domain.ActorBirthDate(time.Now()),
			},
		},
	}

	actorStorage.EXPECT().AreExists(actorIds).Return(true, nil).Times(1)
	filmStorage.EXPECT().Create(filmTitle, filmDescription, filmReleaseDate, filmRating, actorIds).Return(expected, nil).Times(1)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)
	film, err := filmService.Create(filmTitle, filmDescription, filmReleaseDate, filmRating, actorIds)
	require.NoError(t, err)
	require.Equal(t, expected, film)
}

func TestCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)

	type in struct {
		Title       domain.FilmTitle
		Description domain.FilmDescription
		ReleaseDate domain.FilmReleaseDate
		Rating      domain.FilmRating
		ActorIds    []domain.ActorId
	}

	var (
		filmTitle1 = domain.FilmTitle("Title_1")
		filmTitle2 = domain.FilmTitle("Title_2")

		filmDescription1 = domain.FilmDescription("Description_1")
		filmDescription2 = domain.FilmDescription("Description_2")

		filmReleaseDate1 = domain.FilmReleaseDate(time.Now())
		filmReleaseDate2 = domain.FilmReleaseDate(time.Now().Add(1 * time.Minute))

		filmRating1 = domain.FilmRating(5)
		filmRating2 = domain.FilmRating(10)

		actorIds1 = []domain.ActorId{1, 2}
		actorIds2 = []domain.ActorId{1, 2, 3}
	)

	testCases := []struct {
		name         string
		in           in
		out          *domain.Film
		createExpErr error
		exists       bool
		existsExpErr error
	}{
		{
			name: "err_film_actors_not_found",
			in: in{
				Title:       filmTitle1,
				Description: filmDescription1,
				ReleaseDate: filmReleaseDate1,
				Rating:      filmRating1,
				ActorIds:    actorIds1,
			},
			out:          nil,
			createExpErr: domain.ErrFilmActorsNotFound,
			exists:       false,
			existsExpErr: nil,
		},
		{
			name: "err_other",
			in: in{
				Title:       filmTitle2,
				Description: filmDescription2,
				ReleaseDate: filmReleaseDate2,
				Rating:      filmRating2,
				ActorIds:    actorIds2,
			},
			out:          nil,
			createExpErr: fmt.Errorf("failed to create a film: %w", errors.New("database is down")),
			exists:       true,
			existsExpErr: nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actorStorage.EXPECT().AreExists(tCase.in.ActorIds).Return(tCase.exists, tCase.existsExpErr).AnyTimes()
			filmStorage.EXPECT().Create(tCase.in.Title, tCase.in.Description, tCase.in.ReleaseDate, tCase.in.Rating, tCase.in.ActorIds).Return(tCase.out, tCase.createExpErr).AnyTimes()
			film, err := filmService.Create(tCase.in.Title, tCase.in.Description, tCase.in.ReleaseDate, tCase.in.Rating, tCase.in.ActorIds)
			require.Error(t, err)
			require.EqualError(t, tCase.createExpErr, err.Error())
			require.Nil(t, film)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	filmId := domain.FilmId(1)
	filmTitle := domain.FilmTitle("Title_1")
	filmDescription := domain.FilmDescription("Description_1")
	filmReleaseDate := domain.FilmReleaseDate(time.Now())
	filmRating := domain.FilmRating(10)
	actorIds := []domain.ActorId{1, 2}

	expected := &domain.Film{
		Id:          filmId,
		Title:       filmTitle,
		Description: filmDescription,
		ReleaseDate: filmReleaseDate,
		Rating:      filmRating,
		Actors: []*domain.Actor{
			{
				Id:        domain.ActorId(1),
				Name:      domain.ActorName("Actor_1"),
				Sex:       domain.ActorSex(0),
				BirthDate: domain.ActorBirthDate(time.Now()),
			},
			{
				Id:        domain.ActorId(2),
				Name:      domain.ActorName("Actor_2"),
				Sex:       domain.ActorSex(1),
				BirthDate: domain.ActorBirthDate(time.Now()),
			},
		},
	}

	filmStorage.EXPECT().IsExists(filmId).Return(true, nil).Times(1)
	actorStorage.EXPECT().AreExists(actorIds).Return(true, nil).Times(1)
	filmStorage.EXPECT().Update(filmId, &filmTitle, &filmDescription, &filmReleaseDate, &filmRating, &actorIds).Return(expected, nil).Times(1)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)
	film, err := filmService.Update(filmId, &filmTitle, &filmDescription, &filmReleaseDate, &filmRating, &actorIds)
	require.NoError(t, err)
	require.Equal(t, expected, film)
}

func TestUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)

	type in struct {
		Id          domain.FilmId
		Title       *domain.FilmTitle
		Description *domain.FilmDescription
		ReleaseDate *domain.FilmReleaseDate
		Rating      *domain.FilmRating
		ActorIds    *[]domain.ActorId
	}

	var (
		filmId1          = domain.FilmId(1)
		filmTitle1       = domain.FilmTitle("Title_1")
		filmDescription1 = domain.FilmDescription("Description_1")
		filmReleaseDate1 = domain.FilmReleaseDate(time.Now())
		filmRating1      = domain.FilmRating(10)
		actorIds1        = []domain.ActorId{1, 2}

		filmId2          = domain.FilmId(2)
		filmTitle2       = domain.FilmTitle("Title_2")
		filmDescription2 = domain.FilmDescription("Description_2")
		filmReleaseDate2 = domain.FilmReleaseDate(time.Now())
		filmRating2      = domain.FilmRating(6)
		actorIds2        = []domain.ActorId{}

		filmId3          = domain.FilmId(3)
		filmTitle3       = domain.FilmTitle("Title_3")
		filmDescription3 = domain.FilmDescription("Description_3")
		filmReleaseDate3 = domain.FilmReleaseDate(time.Now())
		filmRating3      = domain.FilmRating(9)
		actorIds3        = []domain.ActorId{2, 3, 4}
	)

	testCases := []struct {
		name                  string
		in                    in
		out                   *domain.Film
		updateExpErr          error
		isFilmExists          bool
		isFilmExistsExpErr    error
		areActorsExists       bool
		areActorsExistsExpErr error
	}{
		{
			name: "err_film_not_found",
			in: in{
				Id:          filmId1,
				Title:       &filmTitle1,
				Description: &filmDescription1,
				ReleaseDate: &filmReleaseDate1,
				Rating:      &filmRating1,
				ActorIds:    &actorIds1,
			},
			out:                   nil,
			updateExpErr:          domain.ErrFilmNotFound,
			isFilmExists:          false,
			isFilmExistsExpErr:    nil,
			areActorsExists:       true,
			areActorsExistsExpErr: nil,
		},
		{
			name: "err_film_actors_not_found",
			in: in{
				Id:          filmId2,
				Title:       &filmTitle2,
				Description: &filmDescription2,
				ReleaseDate: &filmReleaseDate2,
				Rating:      &filmRating2,
				ActorIds:    &actorIds2,
			},
			out:                   nil,
			updateExpErr:          domain.ErrFilmActorsNotFound,
			isFilmExists:          true,
			isFilmExistsExpErr:    nil,
			areActorsExists:       false,
			areActorsExistsExpErr: nil,
		},
		{
			name: "err_other",
			in: in{
				Id:          filmId3,
				Title:       &filmTitle3,
				Description: &filmDescription3,
				ReleaseDate: &filmReleaseDate3,
				Rating:      &filmRating3,
				ActorIds:    &actorIds3,
			},
			out:                   nil,
			updateExpErr:          fmt.Errorf("failed to update a film: %w", errors.New("database is down")),
			isFilmExists:          true,
			isFilmExistsExpErr:    nil,
			areActorsExists:       true,
			areActorsExistsExpErr: nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			filmStorage.EXPECT().IsExists(tCase.in.Id).Return(tCase.isFilmExists, tCase.isFilmExistsExpErr).AnyTimes()
			actorStorage.EXPECT().AreExists(*tCase.in.ActorIds).Return(tCase.areActorsExists, tCase.areActorsExistsExpErr).AnyTimes()
			filmStorage.EXPECT().Update(tCase.in.Id, tCase.in.Title, tCase.in.Description, tCase.in.ReleaseDate, tCase.in.Rating, tCase.in.ActorIds).Return(tCase.out, tCase.updateExpErr).AnyTimes()
			film, err := filmService.Update(tCase.in.Id, tCase.in.Title, tCase.in.Description, tCase.in.ReleaseDate, tCase.in.Rating, tCase.in.ActorIds)
			require.Error(t, err)
			require.EqualError(t, tCase.updateExpErr, err.Error())
			require.Nil(t, film)
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmsStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	filmId := domain.FilmId(1)

	filmsStorage.EXPECT().IsExists(filmId).Return(true, nil).Times(1)
	filmsStorage.EXPECT().Delete(filmId).Return(nil).Times(1)

	filmService := domain.NewFilmService(filmsStorage, actorStorage, logsBuilder)
	err := filmService.Delete(filmId)
	require.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)

	type in struct {
		Id domain.FilmId
	}

	var (
		filmId1 = domain.FilmId(1)
		filmId2 = domain.FilmId(2)
	)

	testCases := []struct {
		name         string
		in           in
		deleteExpErr error
		exists       bool
		existsExpErr error
	}{
		{
			name: "err_film_not_found",
			in: in{
				Id: filmId1,
			},
			deleteExpErr: domain.ErrFilmNotFound,
			exists:       false,
			existsExpErr: nil,
		},
		{
			name: "err_other",
			in: in{
				Id: filmId2,
			},
			deleteExpErr: fmt.Errorf("failed to delete a film: %w", errors.New("database is down")),
			exists:       true,
			existsExpErr: nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			filmStorage.EXPECT().IsExists(tCase.in.Id).Return(tCase.exists, tCase.existsExpErr).AnyTimes()
			filmStorage.EXPECT().Delete(tCase.in.Id).Return(tCase.deleteExpErr).AnyTimes()
			err := filmService.Delete(tCase.in.Id)
			require.Error(t, err)
			require.EqualError(t, tCase.deleteExpErr, err.Error())
		})
	}
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	expected := []*domain.Film{
		{
			Id:          domain.FilmId(2),
			Title:       domain.FilmTitle("Title_2"),
			Description: domain.FilmDescription("Description_2"),
			ReleaseDate: domain.FilmReleaseDate(time.Now().Add(2 * time.Minute)),
			Rating:      domain.FilmRating(1),
			Actors: []*domain.Actor{
				{
					Id:        domain.ActorId(1),
					Name:      domain.ActorName("Actor_1"),
					Sex:       domain.ActorSex(0),
					BirthDate: domain.ActorBirthDate(time.Now()),
				},
			},
		},
		{
			Id:          domain.FilmId(1),
			Title:       domain.FilmTitle("Title_1"),
			Description: domain.FilmDescription("Description_1"),
			ReleaseDate: domain.FilmReleaseDate(time.Now()),
			Rating:      domain.FilmRating(10),
			Actors: []*domain.Actor{
				{
					Id:        domain.ActorId(1),
					Name:      domain.ActorName("Actor_1"),
					Sex:       domain.ActorSex(0),
					BirthDate: domain.ActorBirthDate(time.Now()),
				},
				{
					Id:        domain.ActorId(2),
					Name:      domain.ActorName("Actor_2"),
					Sex:       domain.ActorSex(1),
					BirthDate: domain.ActorBirthDate(time.Now()),
				},
			},
		},
	}

	titleOrder := ""
	releaseDateOrder := ""
	ratingOrder := "asc"

	limit := 100
	offset := 0

	filmStorage.EXPECT().ListWithSort(titleOrder, releaseDateOrder, ratingOrder, limit, offset).Return(expected, nil).Times(1)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)
	films, err := filmService.ListWithSort(titleOrder, releaseDateOrder, ratingOrder, limit, offset)
	require.NoError(t, err)
	require.Equal(t, expected, films)
}

func TestListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	fakeError := errors.New("database is down")
	expectedErr := fmt.Errorf("failed to list actors: %w", fakeError)

	titleOrder := ""
	releaseDateOrder := ""
	ratingOrder := "asc"

	limit := 100
	offset := 0

	filmStorage.EXPECT().ListWithSort(titleOrder, releaseDateOrder, ratingOrder, limit, offset).Return(nil, expectedErr).Times(1)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)
	actors, err := filmService.ListWithSort(titleOrder, releaseDateOrder, ratingOrder, limit, offset)
	require.Error(t, err)
	require.EqualError(t, expectedErr, err.Error())
	require.Nil(t, actors)
}

func TestSearchByFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	expected := []*domain.Film{
		{
			Id:          domain.FilmId(1),
			Title:       domain.FilmTitle("Title_1"),
			Description: domain.FilmDescription("Description_1"),
			ReleaseDate: domain.FilmReleaseDate(time.Now()),
			Rating:      domain.FilmRating(10),
			Actors: []*domain.Actor{
				{
					Id:        domain.ActorId(1),
					Name:      domain.ActorName("Actor_1"),
					Sex:       domain.ActorSex(0),
					BirthDate: domain.ActorBirthDate(time.Now()),
				},
				{
					Id:        domain.ActorId(2),
					Name:      domain.ActorName("Actor_2"),
					Sex:       domain.ActorSex(1),
					BirthDate: domain.ActorBirthDate(time.Now()),
				},
			},
		},
		{
			Id:          domain.FilmId(2),
			Title:       domain.FilmTitle("Title_2"),
			Description: domain.FilmDescription("Description_2"),
			ReleaseDate: domain.FilmReleaseDate(time.Now().Add(2 * time.Minute)),
			Rating:      domain.FilmRating(1),
			Actors: []*domain.Actor{
				{
					Id:        domain.ActorId(1),
					Name:      domain.ActorName("Actor_1"),
					Sex:       domain.ActorSex(0),
					BirthDate: domain.ActorBirthDate(time.Now()),
				},
			},
		},
	}

	title := domain.FilmTitle("Title")
	actorName := domain.ActorName("Actor")

	limit := 100
	offset := 0

	filmStorage.EXPECT().SearchByFilters(title, actorName, limit, offset).Return(expected, nil).Times(1)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)
	films, err := filmService.SearchByFilters(title, actorName, limit, offset)
	require.NoError(t, err)
	require.Equal(t, expected, films)
}

func TestSearchByFiltersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filmStorage := mocks.NewMockFilmStorage(ctrl)
	actorStorage := mocks.NewMockActorStorage(ctrl)

	logsBuilder := logs.New(os.Stdout, nil)

	fakeError := errors.New("database is down")
	expectedErr := fmt.Errorf("failed to search films: %w", fakeError)

	title := domain.FilmTitle("Title")
	actorName := domain.ActorName("Actor")

	limit := 100
	offset := 0

	filmStorage.EXPECT().SearchByFilters(title, actorName, limit, offset).Return(nil, expectedErr).Times(1)

	filmService := domain.NewFilmService(filmStorage, actorStorage, logsBuilder)
	films, err := filmService.SearchByFilters(title, actorName, limit, offset)
	require.Error(t, err)
	require.EqualError(t, expectedErr, err.Error())
	require.Nil(t, films)
}
