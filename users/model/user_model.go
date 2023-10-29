package model

import (
	"github.com/jackc/pgx/pgtype"
)

type User struct {
	ID        uint        `json:"id"`
	Password  string      `json:"password"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	BirthDate pgtype.Date `json:"birth_date"`
	About     string      `json:"about"`
	AvatarUrl string      `json:"avatarUrl"`
}
