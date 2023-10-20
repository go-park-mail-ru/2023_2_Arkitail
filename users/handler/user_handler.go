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
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	tokenString := cookie.Value
	user, err := h.usecase.GetUserInfo(tokenString)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}


	response, err := h.CreateUserResponse("error", user)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	h.WriteResponse(w, http.StatusOK, response)
}


func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := h.ParseUserFromJsonBody(r)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	_, err = h.usecase.Login(user.Username, user.Password)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	sessionCookie, err := h.usecase.CreateSessionCookie(user.Username)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: sessionCookie,
	}
	http.SetCookie(w, cookie)

	h.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	tokenString := cookie.Value
	err = h.usecase.CheckAuth(tokenString)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	h.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	const passlen = 8
	user, err := h.ParseUserFromJsonBody(r)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusUnauthorized, body)
		return
	}

	if len(user.Password) < passlen {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	if !h.usecase.IsValidPassword(user.Password) {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	if !h.usecase.IsValidEmail(user.Email) {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	_, err = h.usecase.Signup(user)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusUnauthorized, body)
		return
	}
	_, err = h.usecase.CreateSessionCookie(user.Username)
	if err != nil {
		body, _ := h.CreateResponse(errTokenInvalid.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}
	h.WriteResponse(w, http.StatusNoContent, nil)
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

func (h *UserHandler) CreateResponse(errorMsg string) ([]byte, error) {
	response := model.ErrorResponse{Error: errorMsg}
	responseJson, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseJson, nil
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
