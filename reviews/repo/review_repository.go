package repo

import (
	"database/sql"
	"fmt"
	"strconv"

	"project/reviews/model"
)

type ReviewRepository struct {
	DB *sql.DB
}

func NewPlaceRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{
		DB: db,
	}
}

func (r *ReviewRepository) AddReview(review *model.Review) error {
	err := r.DB.QueryRow(
		`INSERT INTO place ("user_id", "place_id", "content", "rating")
        VALUES ($1, $2, $3, $4)`,
		review.UserId,
		review.PlaceId,
		review.Content,
		review.Rating,
	).Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error adding place in a database: %v", err)
	}
	return nil
}

func (r *ReviewRepository) GetReviewById(id uint) (*model.Review, error) {
	review := &model.Review{}
	err := r.DB.
		QueryRow("SELECT id, user_id, place_id, content, rating FROM review where id = $1", id).
		Scan(&review.ID, &review.UserId, &review.Content, &review.Rating)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) GetReviewsByUserId(userId uint) (map[string]*model.Review, error) {
	reviews := make(map[string]*model.Review)
	rows, err := r.DB.Query("SELECT id, user_id, place_id, content, rating FROM review where user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		review := &model.Review{}
		err = rows.Scan(&review.ID, &review.UserId, &review.Content, &review.Rating)
		if err != nil {
			return nil, err
		}
		reviews[strconv.FormatUint(uint64(review.ID), 10)] = review
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewsByPlaceId(placeId uint) (map[string]*model.Review, error) {
	reviews := make(map[string]*model.Review)
	rows, err := r.DB.Query("SELECT id, user_id, place_id, content, rating FROM review where place_id = $1", placeId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		review := &model.Review{}
		err = rows.Scan(&review.ID, &review.UserId, &review.Content, &review.Rating)
		if err != nil {
			return nil, err
		}
		reviews[strconv.FormatUint(uint64(review.ID), 10)] = review
	}
	return reviews, nil
}

// Читаем при удалении, наверно можно сделать лучше
func (r *ReviewRepository) DeleteReviewsById(id uint) error {
	var reviewId uint
	err := r.DB.
		QueryRow("DELETE from review where id = $1", id).
		Scan(&reviewId)
	return err
}
