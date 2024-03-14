package user

type UserStorage interface {
	FindByEmail(email string) (*User, error)
}
