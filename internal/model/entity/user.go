package entity

import "github.com/nebiros/krss/internal/controller/output"

type User struct {
	UserID   int    `db:"user.user_id"`
	Email    string `db:"user.email"`
	Password string `db:"user.password"`
}

func (u *User) ToUserOutput() output.UserOutput {
	return output.UserOutput{
		UserID: u.UserID,
		Email:  u.Email,
	}
}
