package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserResponse struct {
	Id          uuid.UUID           `json:"id"`
	Username    string              `json:"username"`
	Email       string              `json:"email"`
	Roles       []UserRolesResponse `json:"roles"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	ActivatedAt *time.Time          `json:"activated_at"`
}

type UserRolesResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateNewUserRequest struct {
	Username             string            `json:"username" validate:"required"`
	Email                string            `json:"email" validate:"required,email"`
	Password             string            `json:"password" validate:"required"`
	PasswordConfirmation string            `json:"password_confirmation" validate:"required,eqfield=Password"`
	Roles                []entity.UserType `json:"roles" validate:"required,dive,user_roles_custom_validation"`
}

type AssignRoleToUserRequest struct {
	UserId uuid.UUID       `json:"user_id" validate:"required"`
	Role   entity.UserType `json:"role" validate:"required,user_roles_custom_validation"`
}
