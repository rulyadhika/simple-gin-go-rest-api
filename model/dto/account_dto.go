package dto

import "time"

type ResendActivationTokenRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResendActivationTokenRespone struct {
	RequestTime            time.Time `json:"request_time"`
	NextRequestAvailableAt time.Time `json:"next_request_available_at"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordRespone struct {
	RequestTime            time.Time `json:"request_time"`
	NextRequestAvailableAt time.Time `json:"next_request_available_at"`
}
