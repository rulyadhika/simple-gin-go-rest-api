package userrepository

import "golang.org/x/crypto/bcrypt"

func (u *UserRoles) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}
