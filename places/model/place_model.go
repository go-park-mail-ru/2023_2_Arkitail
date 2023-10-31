package model

type Place struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
	Cost        string   `json:"cost,omitempty"`
	Adress      string   `json:"adress,omitempty"`
	WebSite     string   `json:"web-site,omitempty"`
	Email       string   `json:"email,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	ReviewCount uint     `json:"review_count,omitempty"`
	OpenTime    string   `json:"open_hour,omitempty"`
	CloseTime   string   `json:"close_hour,omitempty"`
	ImageURL    string   `json:"imageUrl,omitempty"`
}
