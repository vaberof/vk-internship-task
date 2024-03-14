package domain

type ActorStorage interface {
	Create(name ActorName, sex ActorSex, birthDate ActorBirthDate) (*Actor, error)
	Update(id ActorId, name *ActorName, sex *ActorSex, birthDate *ActorBirthDate) (*Actor, error)
	Delete(id ActorId) error
	List(limit, offset int) ([]*Actor, error)
	IsExists(id ActorId) (bool, error)
	AreExists(ids []ActorId) (bool, error)
}
