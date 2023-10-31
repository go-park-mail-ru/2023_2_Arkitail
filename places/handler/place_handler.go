package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"project/places/model"
	"project/places/usecase"
	"project/utils"
)

var (
	ErrInvalidJson      = errors.New("Invalid JSON")
	ErrFailedToAddPlace = errors.New("Failed to add place")
)

type PlaceHandler struct {
	usecase usecase.PlaceUseCase
}

func NewPlaceHandler(usecase *usecase.PlaceUseCase) *PlaceHandler {
	return &PlaceHandler{*usecase}
}

func (h *PlaceHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	var place model.Place
	err := json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(ErrInvalidJson.Error()))
		return
	}

	err = h.usecase.AddPlace(&place)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(ErrFailedToAddPlace.Error()))
		return
	}
	utils.WriteResponse(w, http.StatusCreated, nil)
}

func (h *PlaceHandler) GetPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := h.usecase.GetPlaces()
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}
	h.WritePlacesIntoJsonResponse(w, http.StatusOK, places)
}

func (h *PlaceHandler) WritePlacesIntoJsonResponse(w http.ResponseWriter, status int, objects map[string]*model.Place) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(objects)
}
