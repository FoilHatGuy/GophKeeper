package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 2
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 2
	argonKeyLen  uint32 = 32
	argonSaltLen uint32 = 16
)

var (
	ErrArgonVersion = errors.New("wrong argon version")
	ErrHashInvalid  = errors.New("invalid hash")
)

// HashPassword hashes password and returns argon2id hash
func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("password verification: %w", err)
	}
	passwordHash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	// Base64 encode the salt and hashed password
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(passwordHash)
	// Return a string using the standard encoded hash representation
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version,
		argonMemory, argonTime, argonThreads, b64Salt, b64Hash), nil
}

// ComparePasswordHash verifies password identity
func ComparePasswordHash(expectedPassword, providedPassword string) (bool, error) {
	vals := strings.Split(expectedPassword, "$")
	if len(vals) != 6 {
		return false, ErrHashInvalid
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return false, fmt.Errorf("password verification: %w", err)
	}
	if version != argon2.Version {
		return false, ErrArgonVersion
	}

	var mem, iter uint32
	var thread uint8
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &mem, &iter, &thread)
	if err != nil {
		return false, fmt.Errorf("password verification: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false, fmt.Errorf("password verification: %w", err)
	}

	actualPasswordHash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return false, fmt.Errorf("password verification: %w", err)
	}

	providedPasswordHash := argon2.IDKey([]byte(providedPassword), salt, iter, mem, thread,
		uint32(len(actualPasswordHash)))
	return subtle.ConstantTimeCompare(actualPasswordHash, providedPasswordHash) == 1, nil
}
