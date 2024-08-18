package models

type UserID int

type User struct {
	ID       UserID `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserInfo struct {
	Username string `json:"login"`
	Password string `json:"password"`
}
