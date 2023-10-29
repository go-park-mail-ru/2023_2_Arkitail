package model

import "time"

type Review struct {
	ID           uint      `json:"id"`
	UserId       uint      `json:"user_id"`
	PlaceId      uint      `json:"place_id"`
	Content      string    `json:"content"`
	Rating       uint      `json:"rating"`
	CreationDate time.Time `json:"creation_date"`
}
