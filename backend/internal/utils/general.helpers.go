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

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func UnpackDependencyMap(dependencyMap map[string]any, services *[]string) {
	for service, children := range dependencyMap {
		if !Contains(*services, service) {
			*services = append(*services, service)
		}
		if children != nil {
			UnpackDependencyMap(children.(map[string]any), services)
		}
	}
}
