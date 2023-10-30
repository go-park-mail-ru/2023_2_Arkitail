package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"project/trip/model"
)

type UserRepository struct {
	DB *sql.DB
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrWrongPassword = errors.New("wrong password")
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetCleanUserById(id uint) (*model.User, error) {
	user := &model.User{}
	err := r.DB.
		QueryRow(`SELECT id, name, birth_date, about, avatar_url FROM "user" WHERE id = $1`, id).
		Scan(&user.ID, &user.Name, &user.BirthDate, &user.About, &user.AvatarUrl)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) GetUser(email string) (*model.User, error) {
	user := &model.User{}
	err := r.DB.
		QueryRow(`SELECT id, password, name, email, birth_date, about, avatar_url FROM "user" WHERE email = $1`, email).
		Scan(&user.ID, &user.Password, &user.Name, &user.Email, &user.BirthDate, &user.About, &user.AvatarUrl)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) GetPlacesInTripResponse(tripId uint) (map[string]*model.PlaceInTripResponse, error) {
	places := make(map[string]*model.PlaceInTripResponse)
	rows, err := r.DB.Query(`SELECT id, name, description, cost, image_url FROM place join trip_to_place on place.id = trip_to_place.place_id as res where res.trip_id = $1`, tripId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) AddPlacesToTrip(tripId uint, places map[string]*model.PlaceInTripRequest) error {
	for _, place := range places {
		err := r.DB.QueryRow(
			`INSERT INTO trip_to_place ("place_id", "trip_id", "first_date", "last_date")
			VALUES ($1, $2, $3, $4)`,
			place.PlaceId,
			tripId,
			place.FirstDate,
			place.LastDate,
		).Scan()
		if err == sql.ErrNoRows {
			err = nil
		}
		if err != nil {
			return fmt.Errorf("error adding trip in a database: %v", err)
		}
	}
	return nil
}

func (r *UserRepository) AddTrip(trip *model.Trip) error {
	err := r.DB.QueryRow(
		`INSERT INTO trip ("user_id", "publicity", "description")
        VALUES ($1, $2, $3) returning id`,
		trip.UserId,
		trip.Publicity,
		trip.Description,
	).Scan(&trip.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error adding trip in a database: %v", err)
	}
	return nil
}
