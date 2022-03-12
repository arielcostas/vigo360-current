package admin

import (
	"errors"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}

	if errors.Is(err, bcrypt.ErrHashTooShort) {
		logger.Notice("[validatepassword]: unable to verify password: hash is too short")
	}

	return false
}
