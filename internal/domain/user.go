package domain

type User struct {
	UUID  string `db:"uuid"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
