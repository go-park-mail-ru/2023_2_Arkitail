package repo

import (
	"database/sql"
	"fmt"
	"os"

	"project/reviews/model"
)

type ReviewRepository struct {
	DB *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{
		DB: db,
	}
}

func (r *ReviewRepository) AddReview(review *model.Review) error {
	err := r.DB.QueryRow(
		`INSERT INTO review ("user_id", "place_id", "content", "rating")
        VALUES ($1, $2, $3, $4) returning id, DATE_TRUNC('second', creation_date)`,
		review.UserId,
		review.PlaceId,
		review.Content,
		review.Rating,
	).Scan(&review.ID, &review.CreationDate)
	if err != nil {
		return fmt.Errorf("error adding review in a database: %v", err)
	}
	return nil
}

func (r *ReviewRepository) GetReviewById(id uint) (*model.Review, error) {
	review := &model.Review{}
	err := r.DB.
		QueryRow("SELECT id, user_id, place_id, content, rating, DATE_TRUNC('second', creation_date) FROM review where id = $1", id).
		Scan(&review.ID, &review.UserId, &review.PlaceId, &review.Content, &review.Rating, &review.CreationDate)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) GetReviewsByUserId(userId uint) ([]*model.Review, error) {
	reviews := make([]*model.Review, 0)
	rows, err := r.DB.Query("SELECT id, user_id, place_id, content, rating, DATE_TRUNC('second', creation_date) FROM review where user_id = $1 order by creation_date", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		review := &model.Review{}
		err = rows.Scan(&review.ID, &review.UserId, &review.PlaceId, &review.Content, &review.Rating, &review.CreationDate)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) GetReviewsByPlaceId(placeId uint) (*model.ReviewsWithAuthors, error) {
	reviewsWithAuthors := &model.ReviewsWithAuthors{Authors: make(map[string]*model.ReviewAuthor)}
	rows, err := r.DB.Query(`SELECT review.id, user_id, place_id, content, rating,
							DATE_TRUNC('second', review.creation_date), avatar_url, name, "user".id
							FROM review join "user"
							on review.user_id = "user".id
							where place_id = $1
							order by review.creation_date`, placeId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		review := &model.Review{}
		author := &model.ReviewAuthor{}
		var avatarUrl sql.NullString
		err = rows.Scan(&review.ID, &review.UserId, &review.PlaceId, &review.Content, &review.Rating, &review.CreationDate,
			&avatarUrl, &author.Name, &author.ID)
		if err != nil {
			return nil, err
		}
		if avatarUrl.Valid {
			author.Avatar, err = os.ReadFile(avatarUrl.String)
			if err != nil {
				author.Avatar = []byte("")
				err = nil
			}
		}
		reviewsWithAuthors.Reviews = append(reviewsWithAuthors.Reviews, review)

		_, isPresent := reviewsWithAuthors.Authors[author.ID]
		if !isPresent {
			reviewsWithAuthors.Authors[author.ID] = author
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviewsWithAuthors, nil
}

func (r *ReviewRepository) DeleteReviewsById(id uint) error {
	err := r.DB.
		QueryRow("DELETE from review where id = $1", id).
		Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}
