package storage

import (
	"errors"
	"sync"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json: "email"`
}

var errUsernameNotFound = errors.New("user with this name doesnt exist")
var errUsernameTaken = errors.New("user with this name already exists")
var errWrongPassword = errors.New("wrong password")

type AuthStorage struct {
	users sync.Map
	len   uint
}

func NewAuthStorage() *AuthStorage {
	storage := &AuthStorage{}
	user := &User{1, "rvasily", "love", "love123@gmail.com"}
	storage.users.Store("rvasily", user)
	storage.len = 1
	return storage
}

func (storage *AuthStorage) GetUser(username string) (user *User, found bool) {
	value, found := storage.users.Load(username)
	if found {
		user = value.(*User)
	}
	return
}

func (storage *AuthStorage) AddUser(username string, password string, email string) (err error) {
	_, found := storage.GetUser(username)
	if found {
		return errUsernameTaken
	}

	user := User{
		storage.len + 1,
		username,
		password,
		email,
	}
	storage.len++
	storage.users.Store(username, &user)
	return
}

func (storage *AuthStorage) ComparePassword(username string, password string, email string) (err error) {
	user, found := storage.GetUser(username)
	if !found {
		return errUsernameNotFound
	}
	if user.Password != password {
		return errWrongPassword
	}
	return
}
