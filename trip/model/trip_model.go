package model

import (
	"project/places/model"

	"github.com/jackc/pgx/pgtype"
)

type Trip struct {
	ID          uint   `json:"id"`
	UserId      uint   `json:"user_id"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Publicity   string `json:"publicity"`
}

type PlaceInTripRequest struct {
	PlaceId   uint        `json:"place_id"`
	FirstDate pgtype.Date `json:"first_date"`
	LastDate  pgtype.Date `json:"last_date,omitempty"`
}

type TripRequest struct {
	ID          uint                           `json:"id"`
	UserId      uint                           `json:"user_id"`
	Description string                         `json:"description,omitempty"`
	Name        string                         `json:"name"`
	Publicity   string                         `json:"publicity"`
	Places      map[string]*PlaceInTripRequest `json:"places_in_trip"`
}

type PlaceInTripResponse struct {
	Place     *model.Place
	FirstDate pgtype.Date `json:"first_date"`
	LastDate  pgtype.Date `json:"last_date,omitempty"`
}

type TripResponse struct {
	ID          uint                            `json:"id"`
	UserId      uint                            `json:"user_id"`
	Description string                          `json:"description,omitempty"`
	Name        string                          `json:"name"`
	Publicity   string                          `json:"publicity"`
	Places      map[string]*PlaceInTripResponse `json:"places_in_trip"`
}
