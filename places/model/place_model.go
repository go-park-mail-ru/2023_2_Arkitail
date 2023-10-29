package model

type Place struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
	Cost        string   `json:"cost,omitempty"`
	ImageURL    string   `json:"imageUrl,omitempty"`
}
