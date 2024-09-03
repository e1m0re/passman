package model

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
