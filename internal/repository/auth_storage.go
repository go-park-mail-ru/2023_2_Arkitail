package storage

import "errors"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var errUsernameNotFound = errors.New("user with this name doesnt exist")
var errUsernameTaken = errors.New("user with this name already exists")
var errWrongPassword = errors.New("wrong password")

type AuthStorage struct {
	users map[string]*User
}

func NewAuthStorage() *AuthStorage {
	return &AuthStorage{
		users: map[string]*User{
			"rvasily": {1, "rvasily", "love"},
		},
	}
}

func (storage *AuthStorage) GetUser(username string) (user *User, found bool) {
	user, found = storage.users[username]
	return
}

func (storage *AuthStorage) AddUser(username string, password string) (err error) {
	_, found := storage.GetUser(username)
	if found {
		return errUsernameTaken
	}

	user := User{
		uint(len(storage.users) + 1),
		username,
		password,
	}
	storage.users[username] = &user
	return
}

func (storage *AuthStorage) ComparePassword(username string, password string) (err error) {
	user, found := storage.GetUser(username)
	if !found {
		return errUsernameNotFound
	}
	if user.Password != password {
		return errWrongPassword
	}
	return
}
