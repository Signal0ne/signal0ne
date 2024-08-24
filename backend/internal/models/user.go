package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaimsWithUserData struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}
