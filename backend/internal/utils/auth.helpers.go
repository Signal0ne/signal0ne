package utils

import (
	"fmt"
	"signal0ne/cmd/config"
	"signal0ne/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const ACCESS_TOKEN_EXPIRATION_TIME = time.Minute * 10
const REFRESH_TOKEN_EXPIRATION_TIME = time.Hour * 24

func ComparePasswordHashes(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CreateToken(user models.User, tokenType string) (string, error) {
	var cfg = config.GetInstance()
	var expTime time.Duration
	var SECRET_KEY = []byte(cfg.SignalOneSecret)

	if tokenType == "refresh" {
		expTime = REFRESH_TOKEN_EXPIRATION_TIME
	} else if tokenType == "access" {
		expTime = ACCESS_TOKEN_EXPIRATION_TIME
	} else {
		expTime = time.Second * 0
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":      time.Now().Add(expTime).Unix(),
			"id":       user.Id.Hex(),
			"userName": user.Name,
			"role":     user.Role,
		})

	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func VerifyToken(tokenString string) (string, error) {
	var cfg = config.GetInstance()
	var claims = &models.JWTClaimsWithUserData{}
	var SECRET_KEY = []byte(cfg.SignalOneSecret)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Id, nil
}
