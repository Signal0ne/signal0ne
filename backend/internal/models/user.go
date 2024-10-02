package models

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password,omitempty" bson:"password"`
	PhotoUri string             `json:"photoUri" bson:"photoUri"`
	Role     Role               `json:"role" bson:"role"`
}

type JWTClaimsWithUserData struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}
