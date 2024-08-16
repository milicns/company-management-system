package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Hash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return []byte{}, err
	}
	return hash, nil
}

func Matches(hash []byte, plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil
		default:
			return err
		}
	}
	return nil
}
