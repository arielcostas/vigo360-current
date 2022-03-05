package admin

import "golang.org/x/crypto/bcrypt"

func ValidatePassword(password string, hash string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if res != nil {
		println(res.Error())
		return false
	}
	return true
}
