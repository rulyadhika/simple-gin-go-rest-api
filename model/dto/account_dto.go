package dto

type ResendTokenRequest struct {
	Email string `json:"email" validate:"required,email"`
	Type  string `json:"type" validate:"required,resend_token_type_custom_validation"`
}
