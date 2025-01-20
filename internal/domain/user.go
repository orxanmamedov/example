package domain

type User struct {
	UUID  string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
