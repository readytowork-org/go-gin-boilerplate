package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPlainPassword(HashedPassword, PlainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(PlainPassword)); err != nil {
		return false
	}
	return true
}
