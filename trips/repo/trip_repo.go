package repo

import (
	"database/sql"
	"fmt"
	"strconv"

	"project/trips/model"
)

type TripRepository struct {
	DB *sql.DB
}

func NewTripRepository(db *sql.DB) *TripRepository {
	return &TripRepository{DB: db}
}

func (r *TripRepository) DeleteTripById(tripId uint) error {
	err := r.DB.
		QueryRow("DELETE from trip where id = $1", tripId).
		Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}

func (r *TripRepository) GetTripsByUserId(userId uint) (map[string]*model.Trip, error) {
	trips := map[string]*model.Trip{}
	rows, err := r.DB.
		Query(`SELECT id, description, name, is_public from trip where user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		description := sql.NullString{}
		trip := &model.Trip{}
		err = rows.Scan(&trip.ID, &description, &trip.Name, &trip.Publicity)
		trip.Description = description.String
		if err != nil {
			return nil, err
		}
		trips[strconv.FormatUint(uint64(trip.ID), 10)] = trip
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return trips, err
}

func (r *TripRepository) GetTripById(tripId uint) (*model.Trip, error) {
	trip := &model.Trip{}
	description := sql.NullString{}
	err := r.DB.
		QueryRow(`SELECT user_id, description, name, is_public from trip where id = $1`, tripId).
		Scan(&trip.UserId, &description, &trip.Name, &trip.Publicity)
	trip.Description = description.String
	if err != nil {
		return nil, err
	}

	trip.ID = tripId
	return trip, err
}

func (r *TripRepository) GetPlacesInTripResponse(tripId uint) (map[string]*model.PlaceInTripResponse, error) {
	places := make(map[string]*model.PlaceInTripResponse)
	rows, err := r.DB.Query(`SELECT place.id, name, description, cost, image_url,
		(select avg(rating) from review where review.place_id = place.id) as rating,
		adress, open_time, close_time, web_site, email, phone_number,
		(select count(id) from review where place.id = place_id) as review_count,
		res.first_date, res.last_date, res.id
		FROM place join (select * from trip_to_place where trip_to_place.trip_id = $1) as res
		on place.id = res.place_id`, tripId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reviewCount sql.NullFloat64
		placeInTripDb := &model.PlaceInTripDb{}
		err = rows.Scan(&placeInTripDb.Place.ID, &placeInTripDb.Place.Name, &placeInTripDb.Place.Description,
			&placeInTripDb.Place.Cost, &placeInTripDb.Place.ImageURL, &placeInTripDb.Place.Rating,
			&placeInTripDb.Place.Adress, &placeInTripDb.Place.OpenTime, &placeInTripDb.Place.CloseTime,
			&placeInTripDb.Place.WebSite, &placeInTripDb.Place.Email, &placeInTripDb.Place.PhoneNumber,
			&reviewCount, &placeInTripDb.FirstDate, &placeInTripDb.LastDate, &placeInTripDb.ID)

		placeInTrip := model.PlaceInTripResponseFromDb(placeInTripDb)
		placeInTrip.Place.Rating = &reviewCount.Float64
		if err != nil {
			return nil, err
		}
		places[placeInTrip.ID] = placeInTrip
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return places, err
}

func (r *TripRepository) AddPlacesToTrip(tripId uint, places map[string]*model.PlaceInTripRequest) error {
	for _, place := range places {
		err := r.DB.QueryRow(
			`INSERT INTO trip_to_place ("place_id", "trip_id", "first_date", "last_date")
			VALUES ($1, $2, $3, $4)`,
			place.PlaceId,
			tripId,
			place.FirstDate,
			place.LastDate,
		).Scan()
		if err == sql.ErrNoRows {
			err = nil
		}
		if err != nil {
			return fmt.Errorf("error adding trip in a database: %v", err)
		}
	}
	return nil
}

func (r *TripRepository) AddTrip(trip *model.Trip) error {
	err := r.DB.QueryRow(
		`INSERT INTO trip ("user_id", "publicity", "description")
        VALUES ($1, $2, $3) returning id`,
		trip.UserId,
		trip.Publicity,
		trip.Description,
	).Scan(&trip.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error adding trip in a database: %v", err)
	}
	return nil
}

// func (r *TripRepository) UpdatePlacesInTrip(tripId uint, places map[string]*model.PlaceInTripRequest) error {
// 	for _, place := range places {
// 		err := r.DB.QueryRow(
// 			`IF EXISTS (SELECT id FROM trip_to_place WHERE id = $1)
// 			BEGIN
// 				UPDATE trip_to_place
// 				SET
// 				WHERE UniqueColumn = 'Something';
// 			END
// 			ELSE
// 			BEGIN
// 				INSERT INTO Tbl
// 				SELECT ...
// 			END`,
// 			place.PlaceId,
// 			tripId,
// 			place.FirstDate,
// 			place.LastDate,
// 		).Scan()
// 		if err == sql.ErrNoRows {
// 			err = nil
// 		}
// 		if err != nil {
// 			return fmt.Errorf("error adding trip in a database: %v", err)
// 		}
// 	}
// 	return nil
// }

func (r *TripRepository) UpdateTrip(trip *model.Trip) error {
	_, err := r.DB.Exec(
		`UPDATE trip SET "publicity" = $1, "description" = $2 where id = $3`,
		trip.Publicity,
		trip.Description,
		trip.ID,
	)
	return err
}
