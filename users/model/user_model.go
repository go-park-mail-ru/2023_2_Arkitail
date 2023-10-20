package model

type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type GetUserInfoResponse struct {
    Error string `json:"error"`
    User  User   `json:"user"`
}