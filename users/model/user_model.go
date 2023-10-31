package model

import (
	"github.com/jackc/pgx/pgtype"
)

type User struct {
	ID        uint        `json:"id"`
	Password  string      `json:"password,omitempty"`
	Email     string      `json:"email,omitempty"`
	Name      string      `json:"name,omitempty"`
	BirthDate pgtype.Date `json:"birth_date"`
	About     string      `json:"about,omitempty"`
	AvatarUrl string      `json:"avatar_url,omitempty"`
}
