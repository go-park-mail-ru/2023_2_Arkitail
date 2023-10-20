package model

type Place struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Cost        string  `json:"cost"`
	ImageURL    string  `json:"imageUrl"`
}
