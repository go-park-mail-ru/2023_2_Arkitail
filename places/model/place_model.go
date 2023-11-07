package model

import (
	"database/sql"
	"math"
	"strconv"
)

type Place struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Rating      *float64 `json:"rating"`
	Cost        string   `json:"cost"`
	Adress      string   `json:"adress"`
	WebSite     string   `json:"web-site"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phoneNumber"`
	ReviewCount uint     `json:"reviewCount"`
	OpenTime    string   `json:"openHour"`
	CloseTime   string   `json:"closeHour"`
	ImageURL    string   `json:"imageUrl"`
}

type PlaceDb struct {
	ID          uint
	Name        string
	Description string
	Cost        string
	ImageURL    sql.NullString
	OpenTime    sql.NullString
	CloseTime   sql.NullString
	Adress      string
	WebSite     sql.NullString
	Email       sql.NullString
	PhoneNumber sql.NullString
	Rating      sql.NullFloat64
}

func PlaceDbToPlace(placeDb *PlaceDb) *Place {
	place := &Place{ID: strconv.FormatUint(uint64(placeDb.ID), 10),
		Name: placeDb.Name, Description: placeDb.Description, Cost: placeDb.Cost,
		Adress: placeDb.Adress, ImageURL: placeDb.ImageURL.String,
		WebSite: placeDb.WebSite.String, Email: placeDb.Email.String,
		PhoneNumber: placeDb.PhoneNumber.String}

	if placeDb.OpenTime.Valid && placeDb.CloseTime.Valid {
		place.OpenTime = placeDb.OpenTime.String[:5]
		place.CloseTime = placeDb.CloseTime.String[:5]
	}
	if placeDb.Rating.Valid {
		rating := math.Floor(placeDb.Rating.Float64*100) / 100
		place.Rating = &rating
	}
	return place
}
