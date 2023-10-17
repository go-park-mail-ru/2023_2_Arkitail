package repo

import (
	"math/rand"
	"project/places/model"
)

var store = NewMemoryStore()

func GetAllPlaces() ([]model.Place, error) {
	return store.GetAllPlaces()
}

func AddPlace(place model.Place) error {
	return store.AddPlace(place)
}

// func AddPlace(place model.Place) error {
// 	place.ID = generateUniqueID()

// 	places = append(places, place)
// 	return nil
// }

func generateUniqueID() string {
	return "unique_id" + RandomString(8)
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
