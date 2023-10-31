package model

import (
	"project/places/model"
	"strconv"

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

type PlaceInTripDb struct {
	ID        uint
	Place     *model.PlaceDb
	FirstDate pgtype.Date
	LastDate  pgtype.Date
}

type PlaceInTripResponse struct {
	ID        string `json:"id,omitempty"`
	Place     *model.Place
	FirstDate string `json:"first_date"`
	LastDate  string `json:"last_date,omitempty"`
}

type TripResponse struct {
	ID          string                          `json:"id,omitempty"`
	UserId      string                          `json:"user_id,omitempty"`
	Description string                          `json:"description,omitempty"`
	Name        string                          `json:"name,omitempty"`
	Publicity   string                          `json:"publicity,omitempty"`
	Places      map[string]*PlaceInTripResponse `json:"places_in_trip,omitempty"`
}

func TripFromTripRequest(trip *TripRequest) *Trip {
	return &Trip{UserId: trip.UserId, Description: trip.Description,
		Name: trip.Name, Publicity: trip.Publicity}
}

func TripResponseFromTrip(trip *Trip) *TripResponse {
	return &TripResponse{UserId: strconv.FormatUint(uint64(trip.UserId), 10),
		Description: trip.Description, Name: trip.Name, Publicity: trip.Publicity}
}

func PlaceInTripResponseFromDb(tripDb *PlaceInTripDb) *PlaceInTripResponse {
	placeInTrip := &PlaceInTripResponse{Place: model.PlaceDbToPlace(tripDb.Place)}
	placeInTrip.FirstDate = tripDb.FirstDate.Time.Format("2020-10-08")
	placeInTrip.LastDate = tripDb.FirstDate.Time.Format("2020-10-08")
	placeInTrip.ID = strconv.FormatUint(uint64(tripDb.ID), 10)

	return placeInTrip
}
