package user

type UserRole string

func (userRole *UserRole) String() string {
	return string(*userRole)
}

type User struct {
	Id       int64
	Email    string
	Password string
	Role     UserRole
}
