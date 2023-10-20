package repo

import (
	"errors"
	"sync"

	"project/users/model"
)

type UserRepository struct {
	users sync.Map
	len   uint
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrWrongPassword = errors.New("wrong password")
)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUser(username string) (*model.User, error) {
	value, found := r.users.Load(username)
	if found {
		user := value.(*model.User)
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (r *UserRepository) AddUser(user *model.User) error {
	_, err := r.GetUser(user.Username)
	if err != nil {
		return ErrUserExists
	}

	user.ID = r.len + 1
	r.len++
	r.users.Store(user.Username, user)
	return nil
}

func (r *UserRepository) ComparePassword(username, password string) error {
	user, err := r.GetUser(username)
	if err != nil {
		return err
	}
	if user.Password != password {
		return ErrWrongPassword
	}
	return nil
}
