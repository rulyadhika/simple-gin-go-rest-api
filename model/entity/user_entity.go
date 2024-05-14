package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id          uuid.UUID
	Username    string
	Email       string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ActivatedAt time.Time
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}
