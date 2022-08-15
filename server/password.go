package server

import "golang.org/x/crypto/bcrypt"

// HashPassword converts the password into hash values
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

// CheckPasswordHash checks if the password is match
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
