package auth

import (
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func generateHash(password []byte) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func compare(hash, password []byte) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func isValidEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}
