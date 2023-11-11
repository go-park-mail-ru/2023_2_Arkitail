package usecase

import (
	"errors"

	"project/trips/model"
	"project/trips/repo"
)

type TripUsecase struct {
	repo *repo.TripRepository
}

var (
	errCantChangePlaceInTrip = errors.New("not authorized to change a trip")
)

func NewTripUsecase(repo *repo.TripRepository) *TripUsecase {
	return &TripUsecase{
		repo: repo,
	}
}

func (u *TripUsecase) PatchTrip(tripRequest *model.TripRequest) (*model.TripResponse, error) {
	trip := model.TripFromTripRequest(tripRequest)
	err := u.repo.UpdateTrip(trip)
	if err != nil {
		return nil, err
	}

	tripResponse := model.TripResponseFromTrip(trip)
	return tripResponse, err
}

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

func (u *TripUsecase) GetTripsByUserId(userId uint) ([]*model.TripResponse, error) {
	trips, err := u.repo.GetTripsByUserId(userId)
	if err != nil {
		return nil, err
	}

	tripResponses := make([]*model.TripResponse, 0)
	for _, trip := range trips {
		places, err := u.repo.GetPlacesInTripResponse(trip.ID)
		if err != nil {
			return nil, err
		}

		tripResponse := model.TripResponseFromTrip(trip)
		tripResponse.Places = places

		tripResponses = append(tripResponses, tripResponse)
	}
	return tripResponses, err
}

func (u *TripUsecase) PatchPlaceInTrip(placeInTrip *model.PlaceInTripRequest) error {
	err := u.repo.UpdatePlaceInTrip(placeInTrip)
	return err
}

func (u *TripUsecase) DeletePlaceInTripById(placeInTripId uint) error {
	err := u.repo.DeletePlaceInTripById(placeInTripId)
	return err
}

func (u *TripUsecase) AddPlaceInTripById(tripId uint, placeInTrip *model.PlaceInTripRequest) error {
	err := u.repo.AddPlaceToTrip(tripId, placeInTrip)
	return err
}

func (u *TripUsecase) CheckAuthOfPlaceInTrip(userId uint, placeInTripId uint) (bool, error) {
	ownerId, err := u.repo.GetUserIdOfPlaceInTrip(placeInTripId)
	if err != nil {
		return false, err
	}

	if ownerId != userId {
		return false, nil
	}
	return true, nil
}

func (u *TripUsecase) CheckAuthOfTrip(userId uint, tripId uint) (bool, error) {
	ownerId, err := u.repo.GetUserIdOfTrip(tripId)
	if err != nil {
		return false, err
	}

	if ownerId != userId {
		return false, nil
	}
	return true, nil
}

func (u *TripUsecase) GetPlaceInTripById(placeInTripId uint) (*model.PlaceInTripRequest, error) {
	placeInTripRequest, err := u.repo.GetPlaceInTripById(placeInTripId)
	if err != nil {
		return nil, err
	}
	return placeInTripRequest, nil
}
