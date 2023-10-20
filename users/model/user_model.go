package model

type User struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	WebSite   string `json:"web-site"`
	About     string `json:"about"`
	AvatarUrl string `json:"avatarUrl"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}