package repo

import (
	"database/sql"

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
		return err
	}
	return nil
}

func (r *PlaceRepository) GetPlaces() (map[string]*model.Place, error) {
	places := make(map[string]*model.Place)
	rows, err := r.DB.Query("SELECT id, name, description, cost, image_url FROM place")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		place := &model.Place{}
		err = rows.Scan(&place.ID, &place.Name, &place.Description, &place.Cost, &place.ImageURL)
		if err != nil {
			return nil, err
		}
		places[place.ID] = place
	}
	return places, nil
}
