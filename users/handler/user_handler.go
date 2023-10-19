package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"project/users/model"
	"project/users/usecase"
	"time"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: userUsecase}
}

var errTokenInvalid = errors.New("token is invalid")

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateResponse("error", errTokenInvalid.Error()))
		return
	}

	tokenString := cookie.Value
	user, err := h.usecase.GetUserInfo(tokenString)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateResponse("error", err.Error()))
		return
	}

	response, err := h.CreateUserResponse("error", user)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateResponse("error", err.Error()))
		return
	}

	h.WriteResponse(w, http.StatusOK, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := h.ParseUserFromJsonBody(r)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateResponse("error", err.Error()))
		return
	}

	cookie, err := h.usecase.Login(user.Username, user.Password)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateResponse("error", err.Error()))
		return
	}
	http.SetCookie(w, cookie)

	h.WriteResponse(w, http.StatusOK, h.CreateResponse("error", ""))
}

func (h *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateResponse("error", errTokenInvalid.Error()))
		return
	}

	tokenString := cookie.Value
	err = h.usecase.CheckAuth(tokenString)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateResponse("error", err.Error()))
		return
	}

	h.WriteResponse(w, http.StatusOK, h.CreateResponse("error", ""))
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	passlen := 8
	user, err := h.ParseUserFromJsonBody(r)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateResponse("error", err.Error()))
		return
	}

	if len(user.Password) < passlen {
		h.WriteResponse(w, http.StatusBadRequest, h.CreateResponse("error", "Password should be at least 8 characters long"))
		return
	}

	if !h.usecase.IsValidPassword(user.Password) {
		h.WriteResponse(w, http.StatusBadRequest, h.CreateResponse("error", "Password should contain letters, digits, and special characters"))
		return
	}

	if !h.usecase.IsValidEmail(user.Email) {
		h.WriteResponse(w, http.StatusBadRequest, h.CreateResponse("error", "Email should contain @ and letters, digits, or special characters"))
		return
	}

	err = h.usecase.Signup(user)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateResponse("error", err.Error()))
		return
	}
	cookie, err := h.usecase.CreateSessionCookie(user.Username)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateResponse("error", err.Error()))
		return
	}

	http.SetCookie(w, cookie)
	h.WriteResponse(w, http.StatusOK, h.CreateResponse("error", ""))
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
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
}

func (h *UserHandler) WriteResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *UserHandler) ParseUserFromJsonBody(r *http.Request) (*model.User, error) {
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		return nil, usecase.ErrInvalidCredentials
	}
	return &user, nil
}

func (h *UserHandler) CreateResponse(errorMsg string, errMessage string) []byte {
	response := model.GetUserInfoResponse{Error: errorMsg, User: model.User{Username: errMessage}}
	responseJson, err := json.Marshal(response)
	if err != nil {
		return nil
	}
	return responseJson
}

func (h *UserHandler) CreateUserResponse(errorMsg string, user *model.User) ([]byte, error) {
	if user == nil {
		user = &model.User{Username: ""}
	}
	response := model.GetUserInfoResponse{Error: errorMsg, User: *user}
	responseJson, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseJson, nil
}
