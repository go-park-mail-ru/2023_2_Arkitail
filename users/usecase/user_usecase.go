package usecase

import (
	"errors"
	"project/users/model"
	"project/users/repo"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserUseCase interface {
    GetUserInfo(tokenString string) (*model.User, error)
    Login(username, password string) (string, error)
    CheckAuth(tokenString string) error
    Signup(user *model.User) (string, error)
    Logout() error
}

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrInvalidCredentials = errors.New("invalid username or password")
)

type UserClaim struct {
	Username string
	jwt.RegisteredClaims
}

type AuthConfig struct {
    Secret []byte
}

type UserUsecase struct {
    repo   repo.UserRepository
    config AuthConfig
}

func NewUserUsecase(repo *repo.UserRepository, config AuthConfig) *UserUsecase {
    return &UserUsecase{
        repo:   *repo,
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

func (u *UserUsecase) CreateSessionCookie(userName string) (string, error) {
	claim := UserClaim{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	str, err := token.SignedString(u.config.Secret)
	if err != nil {
		return "", err
	}

	return str, nil
}

func (u *UserUsecase) GetUserInfo(tokenString string) (*model.User, error) {
	userClaim, err := u.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := u.repo.GetUser(userClaim.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) Login(username, password string) (string, error) {
	user, err := u.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if user.Password != password {
		return "", ErrInvalidCredentials
	}

	token, err := u.CreateSessionCookie(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) CheckAuth(tokenString string) error {
	_, err := u.ValidateToken(tokenString)
	return err
}

func (u *UserUsecase) Signup(user *model.User) (string, error) {
	err := u.repo.AddUser(user)
	if err != nil {
		return "", err
	}

	token, err := u.CreateSessionCookie(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) Logout() error {
	return nil
}

func (u *UserUsecase) ValidateToken(tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return u.config.Secret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, ErrInvalidToken
	}
}
