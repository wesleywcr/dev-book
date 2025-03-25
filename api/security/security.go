package security

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerificatedPassoword(passowordWithHash, passwordString string) error {
	return bcrypt.CompareHashAndPassword([]byte(passowordWithHash), []byte(passwordString))
}
