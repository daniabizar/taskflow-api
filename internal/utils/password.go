package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword - mengubah password plain text menjadi hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword - cek apakah password yang diinput cocok dengan hash di database
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
