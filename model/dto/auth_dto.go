package dto

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
