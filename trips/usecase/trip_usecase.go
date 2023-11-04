package usecase

import (
	"project/trips/model"
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

// TODO: Patchtrip

func (u *TripUsecase) DeleteTripById(tripId uint) error {
	err := u.repo.DeleteTripById(tripId)
	return err
}

func (u *TripUsecase) AddTrip(tripRequest *model.TripRequest) (*model.TripResponse, error) {
	trip := model.TripFromTripRequest(tripRequest)
	err := u.repo.AddTrip(trip)
	if err != nil {
		return nil, err
	}

	err = u.repo.AddPlacesToTrip(trip.ID, tripRequest.Places)
	if err != nil {
		u.repo.DeleteTripById(trip.ID)
		return nil, err
	}

	tripResponsePlaces, err := u.repo.GetPlacesInTripResponse(trip.ID)
	if err != nil {
		u.repo.DeleteTripById(trip.ID)
		return nil, err
	}

	tripResponse := model.TripResponseFromTrip(trip)
	tripResponse.Places = tripResponsePlaces
	return tripResponse, nil
}

func (u *TripUsecase) GetTripReponseById(tripId uint) (*model.TripResponse, error) {
	trip, err := u.repo.GetTripById(tripId)
	if err != nil {
		return nil, err
	}

	places, err := u.repo.GetPlacesInTripResponse(trip.ID)
	if err != nil {
		return nil, err
	}

	tripResponse := model.TripResponseFromTrip(trip)
	tripResponse.Places = places
	return tripResponse, err
}

func (u *TripUsecase) GetTripsByUserId(userId uint) (map[string]*model.TripResponse, error) {
	trips, err := u.repo.GetTripsByUserId(userId)
	if err != nil {
		return nil, err
	}

	tripResponses := make(map[string]*model.TripResponse)
	for id, trip := range trips {
		places, err := u.repo.GetPlacesInTripResponse(trip.ID)
		if err != nil {
			return nil, err
		}

		tripResponse := model.TripResponseFromTrip(trip)
		tripResponse.Places = places

		tripResponses[id] = tripResponse
	}
	return tripResponses, err
}
