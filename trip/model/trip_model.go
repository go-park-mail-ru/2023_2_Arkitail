package model

import "project/places/model"

type Trip struct {
	ID          uint   `json:"id"`
	UserId      uint   `json:"user_id`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Publicity   string `json:"publicity"`
}

type TripRequest struct {
	ID          uint   `json:"id"`
	UserId      uint   `json:"user_id`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Publicity   string `json:"publicity"`
}

type PlaceInTripRequest struct {
	place *model.Place
	first
}

type TripResponse struct {
	ID          uint   `json:"id"`
	UserId      uint   `json:"user_id`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Publicity   string `json:"publicity"`
}
