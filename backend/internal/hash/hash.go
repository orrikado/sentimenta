package hash

import (
	"log"

	c "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashBytes, err := c.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatalf("не удалось хешировать пароль: %v", err)
	}
	return string(hashBytes)
}

func VerifyPassword(password string, hashed_password string) bool {
	if err := c.CompareHashAndPassword([]byte(hashed_password), []byte(password)); err != nil {
		return false
	}
	return true
}
