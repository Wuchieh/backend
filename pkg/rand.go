package pkg

import (
	"crypto/rand"
	"encoding/base64"
	rand2 "math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString 產生指定長度的隨機字串
func GenerateRandomString(length int) (string, error) {
	stringBase64, err := GenerateRandomStringBase64(length)
	if err != nil {
		return "", err
	}
	return stringBase64[:length], nil
}

// GenerateRandomStringBase64 產生指定最短長度的隨機字串，並使用 base64 編碼
func GenerateRandomStringBase64(minLength int) (string, error) {
	randomBytes := make([]byte, minLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	return randomString, nil
}

// GenerateRandomStringSafe 產生指定長度的隨機字串
func GenerateRandomStringSafe(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand2.Intn(len(charset))]
	}
	return string(result)
}
