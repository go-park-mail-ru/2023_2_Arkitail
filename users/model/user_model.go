package model

import (
	"github.com/jackc/pgx/pgtype"
)

type User struct {
	ID        uint        `json:"id"`
	Password  string      `json:"password"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	BirthDate pgtype.Date `json:"birthDate"`
	About     string      `json:"about"`
	AvatarUrl string      `json:"avatarUrl"`
}

type OldUserSignup struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type UserAvatar struct {
	AvatarUrl string `json:"avatarUrl"`
}
