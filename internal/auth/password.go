package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string, params *argon2id.Params) (string, error) {
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", fmt.Errorf("error creating hash: %w", err)
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("error comparing password: %w", err)
	}
	return match, nil
}
