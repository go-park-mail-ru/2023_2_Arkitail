package model

import "time"

type Review struct {
	ID           uint      `json:"id"`
	UserId       uint      `json:"userId"`
	PlaceId      uint      `json:"placeId"`
	Content      string    `json:"content,omitempty"`
	Rating       uint      `json:"rating"`
	CreationDate time.Time `json:"creationDate"`
}
