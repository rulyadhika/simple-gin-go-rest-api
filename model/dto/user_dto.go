package dto

import "time"

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
