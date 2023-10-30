package usecase

import (
	"project/trips/repo"
)

type TripUsecase struct {
	repo *repo.TripRepository
}

func NewTripUsecase(repo *repo.TripRepository) *TripUsecase {
	return &TripUsecase{
		repo: repo,
	}
}

//
