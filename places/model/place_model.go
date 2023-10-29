package model

type Place struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Rating      float64 `json:"rating"`
	Cost        string  `json:"cost,omitempty"`
	ImageURL    string  `json:"imageUrl,omitempty"`
}
