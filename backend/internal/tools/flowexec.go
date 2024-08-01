package tools

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateWebhookSalt() (string, error) {
	randomBytes := make([]byte, 9)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	salt := base64.URLEncoding.EncodeToString(randomBytes)

	return salt, nil
}
