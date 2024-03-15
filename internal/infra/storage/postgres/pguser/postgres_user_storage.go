package pguser

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/vk-internship-task/internal/service/user"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type PgUserStorage struct {
	db *sqlx.DB
}

func NewPgUserStorage(db *sqlx.DB) *PgUserStorage {
	return &PgUserStorage{db: db}
}

func (s *PgUserStorage) FindByEmail(email string) (*user.User, error) {
	query := `
			SELECT id,
			       email,
			       password,
			       role 
			FROM users
			WHERE email=$1
`
	var user PgUser
	row := s.db.QueryRow(query, email)
	if err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find user by email: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return buildUser(&user), nil
}

func buildUser(postgresUser *PgUser) *user.User {
	return &user.User{
		Id:       postgresUser.Id,
		Email:    postgresUser.Email,
		Password: postgresUser.Password,
		Role:     postgresUser.Role,
	}
}
