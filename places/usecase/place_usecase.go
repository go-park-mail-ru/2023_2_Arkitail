package usecase

import (
	"project/places/model"
	"project/places/repo"
)

func GetPlaces() ([]model.Place, error) {
	places, err := repo.GetAllPlaces()
	if err != nil {
		return nil, err
	}
	return places, nil
}

func AddPlace(place model.Place) error {
	err := repo.AddPlace(place)
	if err != nil {
		return err
	}
	return nil
}
