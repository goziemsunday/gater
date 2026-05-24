package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

type Token struct {
	Plaintext string
	Hash      string
}

func GenerateToken() (*Token, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	plaintext := hex.EncodeToString(b)

	hash := sha256.Sum256([]byte(plaintext))

	return &Token{
		Plaintext: plaintext,
		Hash:      hex.EncodeToString(hash[:]),
	}, nil
}

func CompareToken(token string, hash string) bool {
	hashedToken := sha256.Sum256([]byte(token))

	storedHash, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(hashedToken[:], storedHash) == 1
}
