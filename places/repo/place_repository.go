package repo

import (
	"database/sql"
	"fmt"
	"math"

	"project/places/model"
)

type PlaceRepository struct {
	DB *sql.DB
}

func NewPlaceRepository(db *sql.DB) *PlaceRepository {
	return &PlaceRepository{
		DB: db,
	}
}

func (r *PlaceRepository) AddPlace(place *model.Place) error {
	err := r.DB.QueryRow(
		`INSERT INTO place ("name", "description", "cost", "image_url")
        VALUES ($1, $2, $3, $4)`,
		place.Name,
		place.Description,
		place.Cost,
		place.ImageURL,
	).Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error adding place in a database: %v", err)
	}
	return nil
}

func (r *PlaceRepository) GetPlaces() (map[string]*model.Place, error) {
	places := make(map[string]*model.Place)
	rows, err := r.DB.Query("SELECT id, name, description, cost, image_url, (select avg(rating) from review where review.place_id = place.id) as rating," +
		"adress, open_time, close_time, web_site, email, phone_number, (select count(id) from review where place.id = place_id) as review_count FROM place")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		place := &model.Place{}
		rating := sql.NullFloat64{}
		openTime := sql.NullTime{}
		closeTime := sql.NullTime{}
		err = rows.Scan(&place.ID, &place.Name, &place.Description, &place.Cost, &place.ImageURL, &rating, &place.Adress, &openTime, &closeTime, &place.WebSite, &place.Email, &place.PhoneNumber, &place.ReviewCount)
		if rating.Valid {
			rating.Float64 = math.Floor(rating.Float64*100) / 100
			place.Rating = &rating.Float64
		}
		if openTime.Valid && closeTime.Valid {
			place.OpenTime = openTime.Time.Format("3:4")
			place.CloseTime = closeTime.Time.Format("3:4")
		}
		if err != nil {
			return nil, err
		}
		places[place.ID] = place
	}
	return places, nil
}
