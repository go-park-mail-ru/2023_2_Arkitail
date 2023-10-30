package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"project/trips/model"
	"project/trips/usecase"
	"project/utils"

	"github.com/gorilla/mux"
)

type TripHandler struct {
	usecase *usecase.TripUsecase
}

func NewTripHandler(tripUsecase *usecase.TripUsecase) *TripHandler {
	return &TripHandler{usecase: tripUsecase}
}

var (
	errInvalidUrlParam = errors.New("parameter passed by url has wrong format")
	errTokenInvalid    = errors.New("token is invalid")
)

func (h *TripHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	user, err := h.usecase.GetUserFromClaims(userClaim.(*usecase.UserClaim))
	if err != nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(err.Error()))
		return
	}

	response, err := h.CreateUserResponse(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	utils.WriteResponse(w, http.StatusOK, response)
}

func (h *TripHandler) GetCleanUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil || id < 0 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	user, err := h.usecase.GetCleanUserInfoById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		return
	}

	response, err := h.CreateUserResponse(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	utils.WriteResponse(w, http.StatusOK, response)
}

func (h *TripHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	user, err := h.usecase.GetUserFromClaims(userClaim.(*usecase.UserClaim))
	if err != nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(err.Error()))
		return
	}

	err = h.ParseUserFromJsonBody(user, r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.IsValidUser(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.UpdateUser(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	response, err := h.CreateUserResponse(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	utils.WriteResponse(w, http.StatusOK, response)
}

func (h *TripHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := h.ParseUserFromJsonBody(user, r)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	cookie, err := h.usecase.Login(user.Email, user.Password)
	if err != nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}
	http.SetCookie(w, cookie)

	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *TripHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *TripHandler) Signup(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := h.ParseUserFromJsonBody(user, r)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.IsValidUser(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.Signup(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(err.Error()))
		return
	}

	cookie, err := h.usecase.CreateSessionCookie(user)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	http.SetCookie(w, cookie)
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *TripHandler) Logout(w http.ResponseWriter, r *http.Request) {
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
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *TripHandler) ParseUserFromJsonBody(user *model.User, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		return usecase.ErrInvalidCredentials
	}
	return nil
}

func (h *TripHandler) CreateUserResponse(user *model.User) ([]byte, error) {
	responseJson, err := json.Marshal(user)
	return responseJson, err
}
