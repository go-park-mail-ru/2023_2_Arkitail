package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	storage "project/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

var errWrongJsonFormat = errors.New("json has wrong format")
var errTokenInvalid = errors.New("token is invalid")

type AuthHandler struct {
	secret  []byte
	storage *storage.AuthStorage
}

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type AuthResponse struct {
	Error string `json:"error"`
}

type User struct {
	Login string `json:"login"`
}

type GetUserInfoResponse struct {
	Error string `json:"error"`
	User  User   `json:"user"`
}

type UserResponse struct {
}

type UserClaim struct {
	Username string
	jwt.RegisteredClaims
}

func NewAuthHandler(newSecret string) *AuthHandler {
	storage := storage.NewAuthStorage()
	handler := &AuthHandler{secret: []byte(newSecret), storage: storage}
	return handler
}

func (api *AuthHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		body, _ := CreateUserResponseJson(errTokenInvalid.Error(), nil)
		WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	tokenString := cookie.Value
	user := &UserClaim{}
	token, err := jwt.ParseWithClaims(tokenString, user,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return api.secret, nil
		})
	if err != nil {
		body, _ := CreateUserResponseJson(err.Error(), nil)
		WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	if _, ok := token.Claims.(*UserClaim); ok && token.Valid {
		body, _ := CreateUserResponseJson("", &User{user.Username})
		WriteResponse(w, http.StatusOK, body)
		return
	}
	body, _ := CreateUserResponseJson(errTokenInvalid.Error(), nil)
	WriteResponse(w, http.StatusUnauthorized, body)
	return
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := ParseAuthUserFromJsonBody(r)
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	err = api.storage.ComparePassword(user.Login, user.Password, user.Email)
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	cookie, err := api.createSessionCookie(user.Login)
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	http.SetCookie(w, cookie)
	body, err := CreateAuthResponseJson("")
	WriteResponse(w, http.StatusOK, body)
	return
}

func (api *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		body, _ := CreateAuthResponseJson(errTokenInvalid.Error())
		WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return api.secret, nil
		})
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	if _, ok := token.Claims.(*UserClaim); ok && token.Valid {
		body, _ := CreateAuthResponseJson("")
		WriteResponse(w, http.StatusOK, body)
		return
	}
	body, _ := CreateAuthResponseJson(errTokenInvalid.Error())
	WriteResponse(w, http.StatusUnauthorized, body)
	return
}

func (api *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	user, err := ParseAuthUserFromJsonBody(r)
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	if len(user.Password) < 8 {
		body, _ := CreateAuthResponseJson("Password should be at least 8 characters long")
		WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	if !isValidPassword(user.Password) {
		body, _ := CreateAuthResponseJson("Password should contain letters, digits, and special characters")
		WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	if !isValidEmail(user.Email){
		body, _ := CreateAuthResponseJson("Email should contain @ and letters, digits, or special characters")
		WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	err = api.storage.AddUser(user.Login, user.Password, user.Email)
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	cookie, err := api.createSessionCookie(user.Login)
	if err != nil {
		body, _ := CreateAuthResponseJson(err.Error())
		WriteResponse(w, http.StatusInternalServerError, body)
		return
	}
	http.SetCookie(w, cookie)

	body, err := CreateAuthResponseJson("")
	WriteResponse(w, http.StatusOK, body)
	return
}

func isValidPassword(password string) bool {
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\]`).MatchString(password)

	return hasLetter && hasDigit && hasSpecialChar
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9\-]+\.[a-z]{2,4}$`).MatchString(email)

	return emailRegex
}

func (api *AuthHandler) createSessionCookie(userName string) (cookie *http.Cookie, err error) {
	claim := UserClaim{
		userName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	str, err := token.SignedString(api.secret)
	if err != nil {
		return
	}

	cookie = &http.Cookie{
		Name:  "session_id",
		Value: str,
	}
	return
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().Add(-7 * 24 * time.Hour)
	_, err := r.Cookie("session_id")
	if err != nil {
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "value",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
	return
}

func WriteResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func ParseAuthUserFromJsonBody(r *http.Request) (user AuthUser, err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		err = errWrongJsonFormat
		return
	}
	return
}

func CreateAuthResponseJson(errorMsg string) (responseJson []byte, err error) {
	response := AuthResponse{Error: errorMsg}
	responseJson, err = json.Marshal(response)
	return
}

func CreateUserResponseJson(errorMsg string, user *User) (responseJson []byte, err error) {
	if user == nil {
		user = &User{""}
	}
	response := GetUserInfoResponse{Error: errorMsg, User: *user}
	responseJson, err = json.Marshal(response)
	return
}
