package dto

import (
	"time"

	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserResponse struct {
	Id        uint32              `json:"id"`
	Username  string              `json:"username"`
	Email     string              `json:"email"`
	Roles     []UserRolesResponse `json:"roles"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type UserRolesResponse struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}

type CreateNewUserRequest struct {
	Username             string            `json:"username" validate:"required"`
	Email                string            `json:"email" validate:"required,email"`
	Password             string            `json:"password" validate:"required"`
	PasswordConfirmation string            `json:"password_confirmation" validate:"required,eqfield=Password"`
	Roles                []entity.UserType `json:"roles" validate:"required,dive,user_roles_custom_validation"`
}
