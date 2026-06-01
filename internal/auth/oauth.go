package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomOAuthState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("auth: generate oauth state: %w", err)
	}

	plaintext := hex.EncodeToString(b)
	return plaintext, nil
}
