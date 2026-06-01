package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	errInvalidHash         = errors.New("auth: argon2id hash is not in the correct format")
	errIncompatibleVariant = errors.New("auth: incompatible argon2 variant")
	errIncompatibleVersion = errors.New("auth: incompatible argon2 version")
)

type HashParams struct {
	Memory     uint32 // amount of memory used in kibibytes
	Iterations uint32 // number of iterations over the memory
	Threads    uint8  // number of threads for parallelism (should be between 1 and runtime.NumCPU)
	SaltLength uint32 // length of the random salt (16 bytes is recommended)
	KeyLength  uint32 // length of the generated key (32 bytes is recommended)
}

var DefaultHashParams = &HashParams{
	Memory:     uint32(64 * 1024), // 64 MB
	Iterations: uint32(3),
	Threads:    uint8(4),
	SaltLength: uint32(16),
	KeyLength:  uint32(32),
}

func HashPassword(password string, params *HashParams) (string, error) {
	var hashParams *HashParams
	if params != nil {
		hashParams = params
	} else {
		hashParams = DefaultHashParams
	}

	salt, err := generateRandomBytes(hashParams.SaltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, hashParams.Iterations,
		hashParams.Memory, hashParams.Threads, hashParams.KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, hashParams.Memory, hashParams.Iterations, hashParams.Threads, b64Salt, b64Key)
	return hash, nil
}

func VerifyPassword(password, hash string) (bool, error) {
	params, salt, key, err := decodeHash(hash)
	if err != nil {
		return false, err
	}

	otherKey := argon2.IDKey([]byte(password), salt, params.Iterations,
		params.Memory, params.Threads, params.KeyLength)

	if subtle.ConstantTimeEq(int32(len(key)), int32(len(otherKey))) == 0 {
		return false, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, nil
	}
	return false, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func decodeHash(hash string) (params *HashParams, salt, key []byte, err error) {
	r := strings.NewReader(hash)

	_, err = fmt.Fscanf(r, "$argon2id$")
	if err != nil {
		return nil, nil, nil, errIncompatibleVariant
	}

	var version int
	_, err = fmt.Fscanf(r, "v=%d$", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	params = &HashParams{}
	_, err = fmt.Fscanf(r, "m=%d,t=%d,p=%d$", &params.Memory, &params.Iterations, &params.Threads)
	if err != nil {
		return nil, nil, nil, err
	}

	rest, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, nil, err
	}
	if bytes.ContainsAny(rest, "\r\n") { // base64 decoder ignores these
		return nil, nil, nil, errInvalidHash
	}

	var i int
	if i = bytes.IndexByte(rest, '$'); i == -1 {
		return nil, nil, nil, errInvalidHash
	}

	b64Enc := base64.RawStdEncoding.Strict()

	salt = make([]byte, b64Enc.DecodedLen(i))
	_, err = b64Enc.Decode(salt, rest[:i])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key = make([]byte, b64Enc.DecodedLen(len(rest)-i-1))
	_, err = b64Enc.Decode(key, rest[i+1:])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
