package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        uint32
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}
