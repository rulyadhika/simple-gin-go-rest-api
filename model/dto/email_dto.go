package dto

type SendActivationTokenEmailRequest struct {
	ToEmailAddress, Subject, Username, Token string
}

type SendResetPasswordTokenEmailRequest struct {
	ToEmailAddress, Subject, Username, Token, ExpiredAt string
}
