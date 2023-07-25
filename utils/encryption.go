package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes() ([]byte, error) {
	// hardcoded length of []bytes => 16
	b := make([]byte, 16)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe() (string, error) {
	b, err := generateRandomBytes()
	return base64.URLEncoding.EncodeToString(b), err
}
