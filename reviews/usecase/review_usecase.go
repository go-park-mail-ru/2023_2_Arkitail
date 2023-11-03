package usecase

import (
	"project/reviews/model"
	"project/reviews/repo"
)

type ReviewUseCase struct {
	repo *repo.ReviewRepository
}

func NewReviewUsecase(repo *repo.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{
		repo: repo,
	}
}

func (u *ReviewUseCase) GetReviewById(id uint) (*model.Review, error) {
	review, err := u.repo.GetReviewById(id)
	if err != nil {
		return nil, err
	}

	return review, err
}

func (u *ReviewUseCase) GetReviewsByUserId(userId uint) (map[string]*model.Review, error) {
	reviews, err := u.repo.GetReviewsByUserId(userId)
	if err != nil {
		return nil, err
	}

	return reviews, err
}

func (u *ReviewUseCase) GetReviewsByPlaceId(placeId uint) (map[string]*model.Review, error) {
	reviews, err := u.repo.GetReviewsByPlaceId(placeId)
	if err != nil {
		return nil, err
	}

	return reviews, err
}

func (u *ReviewUseCase) AddReview(review *model.Review) error {
	err := u.repo.AddReview(review)
	if err != nil {
		return err
	}
	return nil
}

func (u *ReviewUseCase) DeleteReviewById(id uint) error {
	err := u.repo.DeleteReviewsById(id)
	return err
}
