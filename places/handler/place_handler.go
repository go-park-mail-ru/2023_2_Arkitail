package handler

import (
	"net/http"
	"encoding/json"
	"project/places/model"
	"project/places/usecase"
)

func CreatePlace(w http.ResponseWriter, r *http.Request) {
	var place model.Place
	err := json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = usecase.AddPlace(place)
	if err != nil {
		http.Error(w, "Failed to add place", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := usecase.GetPlaces()
	if err != nil {
		http.Error(w, "Failed to get places", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(places)
}
