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
		return "", fmt.Errorf("auth.GenerateRandomOAuthState: %w", err)
	}

	plaintext := hex.EncodeToString(b)
	return plaintext, nil
}
