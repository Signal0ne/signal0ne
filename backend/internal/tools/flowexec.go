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

func TraverseOutput(
	payload any,
	desiredKey string,
	mapping string) any {

	switch v := payload.(type) {
	case map[string]any:
		for key, value := range v {
			if key == mapping {
				return value
			}
			TraverseOutput(value, desiredKey, mapping)
		}
	case []any:
		for _, value := range v {
			TraverseOutput(value, desiredKey, mapping)
		}
	default:
		return v
	}
	return nil
}
