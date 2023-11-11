package usecase

import (
	"fmt"

	"project/places/model"
	"project/places/repo"
)

type PlaceUseCase struct {
	repo repo.PlaceRepository
}

func NewPlaceUseCase(repo *repo.PlaceRepository) *PlaceUseCase {
	return &PlaceUseCase{*repo}
}

func (uc *PlaceUseCase) AddPlace(place *model.Place) error {
	err := uc.repo.AddPlace(place)
	if err != nil {
		return err
	}
	return nil
}

func (uc *PlaceUseCase) GetPlaces() ([]*model.Place, error) {
	places, err := uc.repo.GetPlaces()
	if err != nil {
		return nil, fmt.Errorf("error getting places from database: %v", err)
	}
	return places, nil
}

func (uc *PlaceUseCase) GetPlaceById(placeId uint) (*model.Place, error) {
	place, err := uc.repo.GetPlaceById(placeId)
	if err != nil {
		return nil, fmt.Errorf("error getting places from database: %v", err)
	}
	return place, nil
}
