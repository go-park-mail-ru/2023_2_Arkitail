package model

type User struct {
	ID        uint   `json:"id" valid:",optional"`
	Username  string `json:"username" valid:"len>3,len<30"`
	Password  string `json:"password" valid:"len>8"`
	Email     string `json:"email valid:"email, optional"`
	Name      string `json:"name" valid:"len>8, len<30, optional"`
	Location  string `json:"location" valid:"len<30, optional"`
	WebSite   string `json:"web-site" valid:"len<30, optional"`
	About     string `json:"about" valid:"-, optional"`
	AvatarUrl string `json:"avatarUrl" valid:"len>10, optional"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
