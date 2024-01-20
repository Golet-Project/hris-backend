package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func Base64String(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	encodedToken := base64.URLEncoding.EncodeToString(buffer)

	return encodedToken, err
}
