package model

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/nebiros/krss/internal/model/entity"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	SignIn(email, password string) (entity.User, error)
}

type User struct {
	*sqlx.DB
}

func NewUser(dbClient *sqlx.DB) *User {
	return &User{
		DB: dbClient,
	}
}

func (m *User) SignIn(email, password string) (entity.User, error) {
	if len(email) <= 0 {
		return entity.User{}, errors.WithStack(&ErrEmptyArgument{
			Name:  "email",
			Value: email,
		})
	}

	if len(password) <= 0 {
		return entity.User{}, errors.WithStack(&ErrEmptyArgument{
			Name:  "password",
			Value: password,
		})
	}

	q := `select
		users.user_id as "user.user_id",
		users.email as "user.email",
		users.password as "user.password"
		from
		users
		where
		lower(users.email) = lower(?)`

	stmt, err := m.Preparex(m.Rebind(q))
	if err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	defer stmt.Close()

	var u entity.User

	if err := stmt.Get(&u, strings.ToLower(strings.TrimSpace(email))); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	if len(u.Password) <= 0 {
		return entity.User{}, errors.WithStack(sql.ErrNoRows)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	return u, nil
}
