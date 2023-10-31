package usecase

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"project/users/model"
	"project/users/repo"
	"project/utils"

	"github.com/golang-jwt/jwt/v4"
)

type UserUseCase interface {
	GetUserInfo(tokenString string) (*model.User, error)
	Login(username, password string) (*http.Cookie, error)
	CheckAuth(tokenString string) error
	Signup(user *model.User) error
	Logout() error
	ValidateToken(tokenString string) (*utils.UserClaim, error)
	GetUserFromClaims(UserClaim *utils.UserClaim) (*model.User, error)
	IsValidUser(user *model.User) error
	GetUserInfoById(id uint) (*model.User, error)
	UploadAvatar(image []byte, id uint) (string, error)
}

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrInvalidCredentials     = errors.New("invalid username or password")
	ErrInvalidPasswordLength  = errors.New("password should be at least 8 characters long")
	ErrInvalidPasswordSymbols = errors.New("password should contain letters, digits, and special characters")
	ErrInvalidEmail           = errors.New("email should contain @ and letters, digits, or special characters")
)

type AuthConfig struct {
	Secret []byte
}

type UserUsecase struct {
	repo   *repo.UserRepository
	config AuthConfig
}

func NewUserUsecase(repo *repo.UserRepository, config AuthConfig) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		config: config,
	}
}

func (u *UserUsecase) IsValidPassword(password string) bool {
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\]`).MatchString(password)

	return hasLetter && hasDigit && hasSpecialChar
}

func (u *UserUsecase) IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9\-]+\.[a-z]{2,4}$`).MatchString(email)

	return emailRegex
}

func (u *UserUsecase) UpdateUser(user *model.User) error {
	err := u.repo.UpdateUser(user)
	return err
}

func (u *UserUsecase) CreateSessionCookie(user *model.User) (*http.Cookie, error) {
	claim := utils.UserClaim{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	str, err := token.SignedString(u.config.Secret)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    str,
		HttpOnly: true,
	}

	return cookie, nil
}

func (u *UserUsecase) GetUserInfo(tokenString string) (*model.User, error) {
	userClaim, err := u.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := u.GetUserFromClaims(userClaim)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetUserInfoById(id uint) (*model.User, error) {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetCleanUserInfoById(id uint) (*model.User, error) {
	user, err := u.repo.GetCleanUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) Login(email, password string) (*http.Cookie, error) {
	user, err := u.repo.GetUser(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, ErrInvalidCredentials
	}

	cookie, err := u.CreateSessionCookie(user)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func (u *UserUsecase) CheckAuth(tokenString string) error {
	_, err := u.ValidateToken(tokenString)
	return err
}

func (u *UserUsecase) Signup(user *model.User) error {
	err := u.repo.AddUser(user)
	return err
}

func (u *UserUsecase) Logout() error {
	return nil
}

func (u *UserUsecase) GetUserFromClaims(userClaim *utils.UserClaim) (*model.User, error) {
	user, err := u.repo.GetUserById(userClaim.Id)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (u *UserUsecase) ValidateToken(tokenString string) (*utils.UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &utils.UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return u.config.Secret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*utils.UserClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (u *UserUsecase) IsValidUser(user *model.User) error {
	passlen := 8
	if len(user.Password) < passlen {
		err := ErrInvalidPasswordLength
		return err
	}

	if !u.IsValidPassword(user.Password) {
		err := ErrInvalidPasswordSymbols
		return err
	}

	if !u.IsValidEmail(user.Email) {
		err := ErrInvalidEmail
		return err
	}
	return nil
}

func (u *UserUsecase) UploadAvatar(image []byte, id uint) (string, error) {
	const avatarPath = "../../images/avatars/"

	oldAvatar, err := u.repo.GetUserAvatarUrl(id)
	if err != nil {
		return "", err
	}

	filename := avatarPath + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpeg"
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.Write(image)
	if err != nil {
		return "", err
	}

	file.Sync()

	err = u.repo.UpdateUserAvatar(id, filename)
	if err != nil {
		os.Remove(filename)
	} else {
		os.Remove(oldAvatar)
	}

	return filename, nil
}
