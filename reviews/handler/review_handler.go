package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"project/reviews/model"
	"project/reviews/usecase"
	"project/utils"
)

var (
	errInvalidReview = errors.New("review is invalid")
	errTokenInvalid  = errors.New("token is invalid")
)

// DeleteReview, AddReview, GetReview, GetUserReviews, GetPlaceReviews

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

	review := &model.Review{}
	err := h.ParseReviewFromBody(review, r)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, utils.CreateErrorResponse(err.Error()))
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
