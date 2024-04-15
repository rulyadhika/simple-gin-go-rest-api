package dto

type RegisterUserRequest struct {
	Username             string `json:"username" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type RegisterUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type RefreshTokenRequest struct {
	Token string
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
