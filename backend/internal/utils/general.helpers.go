package utils

import (
	"math/rand"
)

func GenerateRandomString() string {
	var length int = 6
	var charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var bytes = make([]byte, length)

	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}

	return string(bytes)
}
