package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"project/trips/model"
	"project/trips/usecase"
	"project/utils"

	"github.com/gorilla/mux"
)

var (
	errInvalidTripRequest = errors.New("invalid trip request")
	errInvalidUrlParam    = errors.New("invalid parameters passed in url")
)

type TripHandler struct {
	usecase *usecase.TripUsecase
}

func NewTripHandler(tripUsecase *usecase.TripUsecase) *TripHandler {
	return &TripHandler{usecase: tripUsecase}
}

func (h *TripHandler) GetTripByTripId(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(utils.ErrTokenInvalid.Error()))
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["tripId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	tripResponse, err := h.usecase.GetTripReponseById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteTripResponse(w, http.StatusOK, tripResponse)
}

func (h *TripHandler) DeleteTripByTripId(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(utils.ErrTokenInvalid.Error()))
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["tripId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	tripResponse, err := h.usecase.GetTripReponseById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusNoContent, nil)
		return
	}
	if tripResponse.UserId != strconv.FormatUint(uint64(userClaim.(*utils.UserClaim).Id), 10) {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(utils.ErrTokenInvalid.Error()))
		return
	}

	err = h.usecase.DeleteTripById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *TripHandler) GetTripsByUserId(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(utils.ErrTokenInvalid.Error()))
		return
	}

	tripResponses, err := h.usecase.GetTripsByUserId(userClaim.(*utils.UserClaim).Id)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteTripResponseMap(w, http.StatusOK, tripResponses)
}

func (h *TripHandler) PostTripsByUserId(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(utils.ErrTokenInvalid.Error()))
		return
	}

	var tripRequest model.TripRequest
	err := ParseTripRequestFromBody(&tripRequest, r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		return
	}

	tripRequest.UserId = userClaim.(*utils.UserClaim).Id
	tripResponse, err := h.usecase.AddTrip(&tripRequest)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteTripResponse(w, http.StatusCreated, tripResponse)
}

func ParseTripRequestFromBody(trip *model.TripRequest, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(trip); err != nil {
		return errInvalidTripRequest
	}
	return nil
}

func (h *TripHandler) WriteTripResponse(w http.ResponseWriter, status int, tripResponse *model.TripResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(tripResponse)
}

func (h *TripHandler) WriteTripResponseMap(w http.ResponseWriter, status int, tripResponses map[string]*model.TripResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(tripResponses)
}
