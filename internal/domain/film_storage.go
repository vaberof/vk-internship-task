package domain

type FilmStorage interface {
	Create(title FilmTitle, description FilmDescription, releaseDate FilmReleaseDate, rating FilmRating, actorIds []ActorId) (*Film, error)
	Update(id FilmId, title *FilmTitle, description *FilmDescription, releaseDate *FilmReleaseDate, rating *FilmRating, actorIds *[]ActorId) (*Film, error)
	Delete(id FilmId) error
	ListWithSort(titleOrder, releaseDateOrder, ratingOrder string, limit, offset int) ([]*Film, error)
	SearchByFilters(title FilmTitle, actorName ActorName, limit, offset int) ([]*Film, error)
	IsExists(id FilmId) (bool, error)
}
