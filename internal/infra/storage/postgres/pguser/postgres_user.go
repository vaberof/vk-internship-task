package pguser

type PgUser struct {
	Id       int64
	Email    string
	Password string
	Role     string
}
