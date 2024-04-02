package handlers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ChechPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TestPassword() {
	password := "secret"
	hash, _ := HashPassword(password)

	fmt.Println("Password:", password)
	fmt.Println("Hash:", hash)

	match := ChechPasswordHash(password, hash)
	fmt.Println("Match:", match)
}
