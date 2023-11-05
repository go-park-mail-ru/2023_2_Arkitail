package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"project/users/model"
	"project/users/usecase"
	"project/utils"

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
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	user, err := h.usecase.GetUserFromClaims(userClaim.(*utils.UserClaim))
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

func (h *UserHandler) GetCleanUser(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	user, err := h.usecase.GetUserFromClaims(userClaim.(*utils.UserClaim))
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	image, err := h.ReadImageFromBody(r)
	if err != nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(err.Error()))
		return
	}

	imageUrl, err := h.usecase.UploadAvatar(image, userClaim.(*utils.UserClaim).Id)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	body, err := h.CreateImageUrlResponse(imageUrl)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	utils.WriteResponse(w, http.StatusCreated, body)
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
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *UserHandler) ParseUserFromJsonBody(user *model.User, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		return usecase.ErrInvalidCredentials
	}
	return nil
}

func (h *UserHandler) CreateUserResponse(user *model.User) ([]byte, error) {
	responseJson, err := json.Marshal(user)
	return responseJson, err
}

func (h *UserHandler) CreateImageUrlResponse(imageUrl string) ([]byte, error) {
	responseJson, err := json.Marshal(imageUrl)
	return responseJson, err
}

func (h *UserHandler) ReadImageFromBody(r *http.Request) ([]byte, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return data, err
	}
	defer r.Body.Close()
	return data, nil
}
