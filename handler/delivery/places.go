package delivery

import (
    "encoding/json"
    "net/http"
    "fmt"
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
