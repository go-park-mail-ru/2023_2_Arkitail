package repo

import (
	"database/sql"
	"errors"

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

func (r *UserRepository) GetUser(username string) (*model.User, error) {
	user := &model.User{}
	err := r.DB.
		QueryRow(`SELECT id, password, name, username, email, location, web_site, about, avatar_url FROM "user" WHERE username = $1`, username).
		Scan(&user.ID, &user.Password, &user.Name, &user.Username, &user.Email, &user.Location, &user.WebSite, &user.About, &user.AvatarUrl)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) GetUserById(id int) (*model.User, error) {
	user := &model.User{}
	err := r.DB.
		QueryRow(`SELECT id, password, name, username, email, location, web_site, about, avatar_url FROM "user" WHERE id = $1`, id).
		Scan(&user.ID, &user.Password, &user.Name, &user.Username, &user.Email, &user.Location, &user.WebSite, &user.About, &user.AvatarUrl)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	_, err := r.DB.Exec(
		`UPDATE items SET "password" = $1`+
			`,"name" = $2`+
			`,"username" = $3`+
			`,"email" = $5`+
			`,"location" = $6`+
			`,"web_site" = $7`+
			`,"about" = $8`+
			`,"avatar_url" = $9`+
			`WHERE id = $10`,
		user.Password,
		user.Username,
		user.Email,
		user.Location,
		user.WebSite,
		user.About,
		user.AvatarUrl,
		user.ID,
	)
	return err
}

func (r *UserRepository) AddUser(user *model.User) error {
	_, err := r.GetUser(user.Username)
	if err == nil {
		return ErrUserExists
	}

	err = r.DB.QueryRow(
		`INSERT INTO "user" ("name", "username", "password", "email", "location", "web_site", "about", "avatar_url")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Name,
		user.Username,
		user.Password,
		user.Email,
		user.Location,
		user.WebSite,
		user.About,
		user.AvatarUrl,
	).Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		//TODO: чужая ошибка, надо бы нормально обрабатывать
		return err
	}
	return nil
}
