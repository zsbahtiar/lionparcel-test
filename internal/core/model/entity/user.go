package entity

type User struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	IsAdmin      bool   `db:"is_admin"`
}
