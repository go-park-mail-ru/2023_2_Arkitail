package model

import (
	"time"
)

type Review struct {
	ID           uint      `json:"id"`
	UserId       uint      `json:"userId"`
	PlaceId      uint      `json:"placeId"`
	Content      string    `json:"content"`
	Rating       uint      `json:"rating"`
	CreationDate time.Time `json:"creationDate"`
}

type ReviewsWithAuthors struct {
	Reviews []*Review
	Authors map[string]*ReviewAuthor
}

type ReviewAuthor struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar []byte `json:"avatar"`
}
