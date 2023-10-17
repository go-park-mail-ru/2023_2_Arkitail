package repo

import (
	"sync"
	"project/places/model"
)

type MemoryStore struct {
	places map[string]model.Place
	mu     sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		places: make(map[string]model.Place),
	}
}

func (ms *MemoryStore) GetAllPlaces() ([]model.Place, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []model.Place
	for _, place := range ms.places {
		result = append(result, place)
	}
	return result, nil
}

func (ms *MemoryStore) AddPlace(place model.Place) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	place.ID = generateUniqueID()
	ms.places[place.ID] = place

	return nil
}
