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
	UserId      uint   `json:"userId"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Publicity   string `json:"publicity"`
}

type PlaceInTripRequest struct {
	ID        uint           `json:"id"`
	PlaceId   uint           `json:"placeId"`
	FirstDate utils.JsonDate `json:"firstDate"`
	LastDate  utils.JsonDate `json:"lastDate"`
}

type TripRequest struct {
	ID          uint                          `json:"id"`
	UserId      uint                          `json:"userId"`
	Description string                        `json:"description"`
	Name        string                        `json:"name"`
	Publicity   string                        `json:"publicity"`
	Places      map[string]PlaceInTripRequest `json:"placesInTrip"`
}

type PlaceInTripDb struct {
	ID        uint
	Place     model.PlaceDb
	FirstDate pgtype.Date
	LastDate  pgtype.Date
}

type PlaceInTripResponse struct {
	ID        string      `json:"id"`
	Place     model.Place `json:"place"`
	FirstDate string      `json:"firstDate"`
	LastDate  string      `json:"lastDate"`
}

type TripResponse struct {
	ID          string                         `json:"id"`
	UserId      string                         `json:"userId"`
	Description string                         `json:"description"`
	Name        string                         `json:"name"`
	Publicity   string                         `json:"publicity"`
	Places      map[string]PlaceInTripResponse `json:"placesInTrip"`
}

func TripFromTripRequest(trip *TripRequest) *Trip {
	return &Trip{ID: trip.ID, UserId: trip.UserId, Description: trip.Description,
		Name: trip.Name, Publicity: trip.Publicity}
}

func TripResponseFromTrip(trip *Trip) *TripResponse {
	return &TripResponse{ID: strconv.FormatUint(uint64(trip.ID), 10), UserId: strconv.FormatUint(uint64(trip.UserId), 10),
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
	} else {
		tripBd.Publicity = trip.Publicity
	}
	if trip.Description != "" {
		tripBd.Description.String = trip.Description
		tripBd.Description.Valid = true
	}
	return tripBd
}
