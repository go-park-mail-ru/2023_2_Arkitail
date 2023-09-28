package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Place struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Cost        string  `json:"cost"`
	ImageURL    string  `json:"imageUrl"`
}

var places = map[string]Place{}

func createPlace(w http.ResponseWriter, r *http.Request) {
	var newPlace Place
	err := json.NewDecoder(r.Body).Decode(&newPlace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newPlace.ID = "some-unique-id"
	places[newPlace.ID] = newPlace
	w.WriteHeader(http.StatusCreated)
}

func getPlaces(w http.ResponseWriter, r *http.Request) {
	placeList := []Place{}
	for _, place := range places {
		placeList = append(placeList, place)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(placeList)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Привет, мир!")
	})
	
	r.HandleFunc("/places", createPlace).Methods("POST")
	r.HandleFunc("/places", getPlaces).Methods("GET")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
