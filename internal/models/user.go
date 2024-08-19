package models

type UserID int

type User struct {
	ID       UserID `db:"id"`
	Username []byte `db:"username"`
	Password []byte `db:"password"`
}

type Credentials struct {
	Password []byte `json:"password"`
	Username []byte `json:"username"`
}
