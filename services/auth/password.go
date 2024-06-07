package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword function to hash the given password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword function to compare the hashed password with the user passwor
func ComparePassword(hashPassword string, userPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(userPassword)); err != nil {
		return false
	}

	return true
}