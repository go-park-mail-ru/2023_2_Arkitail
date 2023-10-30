package utils

import "github.com/golang-jwt/jwt/v4"

type UserClaim struct {
	Id uint
	jwt.RegisteredClaims
}
