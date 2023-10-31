package model

import (
	"database/sql"
	"math"
	"strconv"
)

type Place struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
	Cost        string   `json:"cost,omitempty"`
	Adress      string   `json:"adress,omitempty"`
	WebSite     string   `json:"web-site,omitempty"`
	Email       string   `json:"email,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	ReviewCount uint     `json:"review_count,omitempty"`
	OpenTime    string   `json:"open_hour,omitempty"`
	CloseTime   string   `json:"close_hour,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
}

// Не набирает поля для агрегации(review count)
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
