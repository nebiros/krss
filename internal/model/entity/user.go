package entity

type User struct {
	UserID   int    `db:"user_id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
