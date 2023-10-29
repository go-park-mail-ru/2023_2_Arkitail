package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"project/users/model"
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

func (r *UserRepository) GetUserById(id uint) (*model.User, error) {
	user := &model.User{}
	err := r.DB.
		QueryRow(`SELECT id, password, name, email, birth_date, about, avatar_url FROM "user" WHERE id = $1`, id).
		Scan(&user.ID, &user.Password, &user.Name, &user.Email, &user.BirthDate, &user.About, &user.AvatarUrl)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	_, err := r.DB.Exec(
		`UPDATE "user" SET "password" = $1`+
			`,"name" = $2`+
			`,"email" = $3`+
			`,"birth_date" = $4`+
			`,"about" = $5`+
			`,"avatar_url" = $6`+
			`WHERE id = $7`,
		user.Password,
		user.Name,
		user.Email,
		user.BirthDate.Time,
		user.About,
		user.AvatarUrl,
		user.ID,
	)
	return err
}

func (r *UserRepository) AddUser(user *model.User) error {
	_, err := r.GetUser(user.Email)
	if err == nil {
		return ErrUserExists
	}

	err = r.DB.QueryRow(
		`INSERT INTO "user" ("name", "password", "email", "birth_date", "about", "avatar_url")
        VALUES ($1, $2, $3, $4, $5, $6)`,
		user.Name,
		user.Password,
		user.Email,
		user.BirthDate.Time,
		user.About,
		user.AvatarUrl,
	).Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error adding user in a database: %v", err)
	}
	return nil
}
