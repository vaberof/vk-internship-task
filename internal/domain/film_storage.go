package domain

type FilmStorage interface {
	Create(title FilmTitle, description FilmDescription, releaseDate FilmReleaseDate, rating FilmRating, actorIds []ActorId) (*Film, error)
	Update(id FilmId, title *FilmTitle, description *FilmDescription, releaseDate *FilmReleaseDate, rating *FilmRating) (*Film, error)
	Delete(id FilmId) error
	ListWithSort(title *FilmTitle, releaseDate *FilmReleaseDate, rating *FilmRating, limit, offset int) ([]*Film, error)
	SearchByFilters(title FilmTitle, actorName ActorName) ([]*Film, error)
	IsExists(id FilmId) (bool, error)
}
