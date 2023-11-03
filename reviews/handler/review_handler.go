package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"project/reviews/model"
	"project/reviews/usecase"
	"project/utils"

	"github.com/gorilla/mux"
)

var (
	errInvalidReview   = errors.New("review is invalid")
	errTokenInvalid    = errors.New("token is invalid")
	errInvalidUrlParam = errors.New("invalid parameters passed in url")
)

type ReviewHandler struct {
	usecase *usecase.ReviewUseCase
}

func NewReviewHandler(usecase *usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{usecase}
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	review := &model.Review{UserId: userClaim.(*utils.UserClaim).Id}
	err := h.ParseReviewFromBody(review, r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		return
	}

	err = h.usecase.AddReview(review)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	response, err := h.CreateReviewResponse(review)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}
	utils.WriteResponse(w, http.StatusCreated, response)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userClaim := r.Context().Value("userClaim")
	if userClaim == nil {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["reviewId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	review, err := h.usecase.GetReviewById(uint(id))
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}
	if review.ID != userClaim.(*utils.UserClaim).Id {
		utils.WriteResponse(w, http.StatusUnauthorized, utils.CreateErrorResponse(errTokenInvalid.Error()))
		return
	}

	err = h.usecase.DeleteReviewById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["reviewId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	review, err := h.usecase.GetReviewById(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	response, err := h.CreateReviewResponse(review)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}
	utils.WriteResponse(w, http.StatusOK, response)
}

func (h *ReviewHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["userId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	reviews, err := h.usecase.GetReviewsByUserId(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteReviewMapResponse(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) GetPlaceReviews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["placeId"])
	if err != nil || id < 1 {
		utils.WriteResponse(w, http.StatusBadRequest, utils.CreateErrorResponse(errInvalidUrlParam.Error()))
		return
	}

	reviews, err := h.usecase.GetReviewsByPlaceId(uint(id))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
		return
	}

	h.WriteReviewMapResponse(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) ParseReviewFromBody(review *model.Review, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(review); err != nil {
		return errInvalidReview
	}
	return nil
}

func (h *ReviewHandler) CreateReviewResponse(review *model.Review) ([]byte, error) {
	responseJson, err := json.Marshal(review)
	return responseJson, err
}

func (h *ReviewHandler) WriteReviewMapResponse(w http.ResponseWriter, status int, reviews map[string]*model.Review) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(reviews)
}
