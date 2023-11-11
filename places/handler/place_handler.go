package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"project/places/model"
	"project/places/usecase"
	"project/utils"

	"github.com/gorilla/mux"
)

var (
	ErrInvalidJson      = errors.New("Invalid JSON")
	ErrFailedToAddPlace = errors.New("Failed to add place")
	errInvalidUrlParam  = errors.New("invalid parameters passed in url")
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

func (h *PlaceHandler) GetPlace(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["placeId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	place, err := h.usecase.GetPlaceById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}
	h.WritePlaceIntoJsonResponse(w, http.StatusOK, place)
}

func (h *PlaceHandler) WritePlacesIntoJsonResponse(w http.ResponseWriter, status int, objects []*model.Place) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(objects)
}

func (h *PlaceHandler) WritePlaceIntoJsonResponse(w http.ResponseWriter, status int, place *model.Place) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(place)
}
