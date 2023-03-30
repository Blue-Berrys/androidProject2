package Bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptionByPassword(password string) string {
	password2 := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(password2, bcrypt.DefaultCost)
	return string(hashedPassword)
}
