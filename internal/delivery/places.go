package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Place struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Cost        string  `json:"cost"`
	ImageURL    string  `json:"imageUrl"`
}

var places = map[string]Place{
	"1": {
		ID:          "1",
		Name:        "Эфелева башня",
		Description: "Это знаменитое архитектурное сооружение, которое находится в центре Парижа, Франция. Эта башня является одной из самых узнаваемых и посещаемых достопримечательностей мира, а также символом как самого Парижа, так и Франции в целом. Она была построена для Всемирной выставки 1889 года, которая отмечала столетие Великой французской революции.",
		Rating:      4.5,
		Cost:        "$$",
		ImageURL:    "https://example.com/image1.jpg",
	},
	"2": {
		ID:          "2",
		Name:        "Эрмитаж",
		Description: "Это один из самых знаменитых и крупнейших музеев мира, расположенный в Санкт-Петербурге, Россия. Этот музей является одной из наиболее значимых культурных достопримечательностей России и мировым центром искусства и культуры.",
		Rating:      3.8,
		Cost:        "$",
		ImageURL:    "https://example.com/image2.jpg",
	},
	"3": {
		ID:          "3",
		Name:        "МГТУ им. Баумана",
		Description: "является одним из ведущих технических университетов в России и весьма престижным учебным заведением.",
		Rating:      5.0,
		Cost:        "$$$",
		ImageURL:    "https://example.com/image2.jpg",
	},
}

func CreatePlace(w http.ResponseWriter, r *http.Request) {
	var newPlace Place
	err := json.NewDecoder(r.Body).Decode(&newPlace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newPlace.ID = generateUniqueID()

	places[newPlace.ID] = newPlace

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": newPlace.ID})
}

func GetPlaces(w http.ResponseWriter, r *http.Request) {
	placeList := []Place{}
	for _, place := range places {
		placeList = append(placeList, place)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(placeList)
}

func generateUniqueID() string {
	return fmt.Sprintf("id%d", len(places)+1)
}
