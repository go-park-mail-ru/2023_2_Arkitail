package handler

import (
	"encoding/json"
	"net/http"
	
	"project/places/model"
	"project/places/usecase"
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
  
	err = h.usecase.AddPlace(&place)
	if err != nil {
		http.Error(w, "Failed to add place", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *PlaceHandler) GetPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := h.usecase.GetPlaces()
	if err != nil {
		http.Error(w, "Failed to get places", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(places)
}
