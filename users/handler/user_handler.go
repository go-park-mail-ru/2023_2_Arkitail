package handler

import (
	"encoding/json"
	"errors"
	"log"
  "strconv"
	"net/http"
	"time"

	"project/users/model"
	"project/users/usecase"
  
  "github.com/gorilla/mux"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: userUsecase}
}

var (
	errInvalidUrlParam = errors.New("parameter passed by url has wrong format")
	errTokenInvalid    = errors.New("token is invalid")
)

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	tokenString := cookie.Value
	user, err := h.usecase.GetUserInfo(tokenString)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateErrorResponse(err.Error()))
		return
	}

	response, err := h.CreateUserResponse(user)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteResponse(w, http.StatusOK, response)
}

func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		body, _ := h.CreateErrorResponse(errInvalidUrlParam.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	user, err := h.usecase.GetUserInfoById(id)
	if err != nil {
		body, _ := h.CreateErrorResponse(err.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	err = h.ParseUserFromJsonBody(user, r)
	if err != nil {
		body, _ := h.CreateErrorResponse(err.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	err = h.usecase.IsValidUser(user)
	if err != nil {
		body, _ := h.CreateErrorResponse(err.Error())
		h.WriteResponse(w, http.StatusBadRequest, body)
		return
	}

	err = h.usecase.UpdateUser(user)
	if err != nil {
		body, _ := h.CreateErrorResponse(err.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	response, err := h.CreateUserResponse(user)
	if err != nil {
		body, _ := h.CreateErrorResponse(err.Error())
		h.WriteResponse(w, http.StatusInternalServerError, body)
		return
	}

	h.WriteResponse(w, http.StatusOK, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := h.ParseUserFromJsonBody(user, r)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	cookie, err := h.usecase.Login(user.Username, user.Password)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}
	http.SetCookie(w, cookie)

	h.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	tokenString := cookie.Value
	err = h.usecase.CheckAuth(tokenString)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := h.ParseUserFromJsonBody(user, r)
	if err != nil {
		body, _ := h.CreateErrorResponse(err.Error())
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.IsValidUser(user)
	if err != nil {
		h.WriteResponse(w, http.StatusBadRequest, h.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.Signup(user)
	if err != nil {
		h.WriteResponse(w, http.StatusUnauthorized, h.CreateErrorResponse(err.Error()))
		return
	}

	cookie, err := h.usecase.CreateSessionCookie(user.Username)
	if err != nil {
		h.WriteResponse(w, http.StatusInternalServerError, h.CreateErrorResponse(err.Error()))
		return
	}

	http.SetCookie(w, cookie)
	h.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	if err != nil {
		return
	}

	expire := time.Now().Add(-7 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "value",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func (h *UserHandler) WriteResponse(w http.ResponseWriter, status int, body []byte) {
	if body == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *UserHandler) ParseUserFromJsonBody(user *model.User, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		return usecase.ErrInvalidCredentials
	}
	return nil
}

func (h *UserHandler) CreateErrorResponse(errorMsg string) []byte {
	response := model.ErrorResponse{Error: errorMsg}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return nil
	}
	return responseJson
}

func (h *UserHandler) CreateUserResponse(user *model.User) ([]byte, error) {
	responseJson, err := json.Marshal(user)
	return responseJson, err
}
