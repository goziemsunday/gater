package auth

import (
	"crypto/rand"
	"crypto/sha256"
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

// func CompareToken(token string, hash string) bool {
//     h := sha256.Sum256([]byte(token))

//     storedHash, err := hex.DecodeString(hash)
//     if err != nil {
//         return false
//     }

//     return subtle.ConstantTimeCompare(h[:], storedHash) == 1
// }
// Or compare both as hex strings:
// func CompareToken(token string, hash string) bool {
//     h := sha256.Sum256([]byte(token))
//     return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(h[:])), []byte(hash)) == 1
// }

// func CompareToken(token string, hash string) bool {
// 	h := sha256.Sum256([]byte(token))

// 	storedHash, err := hex.DecodeString(hash)
// 	if err != nil {
// 		return false
// 	}

// 	return subtle.ConstantTimeCompare(h[:], []byte(hash)) == 1
// }
