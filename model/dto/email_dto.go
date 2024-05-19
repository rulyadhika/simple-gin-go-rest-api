package dto

type SendTokenEmailRequest struct {
	ToEmailAddress, Subject, Username, Token string
}
