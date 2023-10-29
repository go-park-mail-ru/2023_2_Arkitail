package model

import "github.com/jackc/pgx/pgtype"

type Review struct {
	ID           uint        `json:"id"`
	UserId       uint        `json:"user_id"`
	PlaceId      uint        `json:"place_id"`
	Content      string      `json:"content"`
	Rating       uint        `json:"rating"`
	CreationDate pgtype.Date `json:"creation_date"`
}
