package domain

type User struct {
	ID       int     `db:"id"`       // Maps to the `ID` column
	Username *string `db:"username"` // Maps to the `Username` column
	Email    *string `db:"email"`    // Maps to the `Email` column
}
