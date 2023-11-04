package model

import (
	"database/sql"
	"project/places/model"
	"project/utils"
	"strconv"
	"time"

	"github.com/jackc/pgx/pgtype"
)

type TripBd struct {
	ID          uint
	UserId      uint
	Description sql.NullString
	Name        string
	Publicity   string
}

type Trip struct {
	ID          uint   `json:"id"`
	UserId      uint   `json:"user_id"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Publicity   string `json:"publicity"`
}

type PlaceInTripRequest struct {
	PlaceId   uint           `json:"place_id"`
	FirstDate utils.JsonDate `json:"first_date,omitempty"`
	LastDate  utils.JsonDate `json:"last_date,omitempty"`
}

type TripRequest struct {
	ID          uint                          `json:"id"`
	UserId      uint                          `json:"user_id"`
	Description string                        `json:"description,omitempty"`
	Name        string                        `json:"name"`
	Publicity   string                        `json:"publicity"`
	Places      map[string]PlaceInTripRequest `json:"places_in_trip"`
}

type PlaceInTripDb struct {
	ID        uint
	Place     model.PlaceDb
	FirstDate pgtype.Date
	LastDate  pgtype.Date
}

type PlaceInTripResponse struct {
	ID        string      `json:"id,omitempty"`
	Place     model.Place `json:"place,omitempty"`
	FirstDate string      `json:"first_date,omitempty"`
	LastDate  string      `json:"last_date,omitempty"`
}

type TripResponse struct {
	ID          string                         `json:"id,omitempty"`
	UserId      string                         `json:"user_id,omitempty"`
	Description string                         `json:"description,omitempty"`
	Name        string                         `json:"name,omitempty"`
	Publicity   string                         `json:"publicity,omitempty"`
	Places      map[string]PlaceInTripResponse `json:"places_in_trip,omitempty"`
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
	placeInTrip := &PlaceInTripResponse{Place: *model.PlaceDbToPlace(&tripDb.Place)}
	if tripDb.FirstDate.Status == pgtype.Present {
		placeInTrip.FirstDate = tripDb.FirstDate.Time.Format(time.DateOnly)
		if tripDb.LastDate.Status == pgtype.Present {
			placeInTrip.LastDate = tripDb.FirstDate.Time.Format(time.DateOnly)
		}
	}
	placeInTrip.ID = strconv.FormatUint(uint64(tripDb.ID), 10)

	return placeInTrip
}

func TripToTripBd(trip *Trip) *TripBd {
	tripBd := &TripBd{ID: trip.ID, UserId: trip.UserId, Name: trip.Name}
	if trip.Publicity == "" {
		tripBd.Publicity = "private"
	}
	if trip.Description != "" {
		tripBd.Description.String = trip.Description
		tripBd.Description.Valid = true
	}
	return tripBd
}
