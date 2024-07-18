package helper

import (
	"fmt"
	"log"

	"github.com/matcornic/hermes/v2"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
)

type EmailHTML interface {
	GenerateSendActivationTokenEmailHTML(username string, token string) (*string, error)
	GenerateSendPasswordResetTokenEmailHTML(username string, token string, expiredAt string) (*string, error)
}

type emailHTMLImpl struct{}

func NewEmailHTMLImpl() EmailHTML {
	return &emailHTMLImpl{}
}

func (e *emailHTMLImpl) GenerateSendActivationTokenEmailHTML(username string, token string) (*string, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				fmt.Sprintf("Welcome to IT Helpdesk System For %s! We're very excited to have you on board.", config.GetAppConfig().EMAIL_SENDER_IDENTITY),
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started, please activate your account by clicking here:",
					Button: hermes.Button{
						Color: "#196ECD",
						Text:  "Activate your account",
						Link:  fmt.Sprintf("%s/%s", config.GetAppConfig().ACCOUNT_ACTIVATION_URL, token),
					},
				},
			},
			Outros: []string{
				"This is an automated email, please don't send or reply to this email.",
			},
		},
	}

	htmlString, err := e.generateHTML(&email)

	if err != nil {
		return nil, err
	}

	return htmlString, nil
}

func (e *emailHTMLImpl) GenerateSendPasswordResetTokenEmailHTML(username string, token string, expiredAt string) (*string, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				fmt.Sprintf("You have received this email because a password reset request for %s IT Helpdesk System account was received.", config.GetAppConfig().EMAIL_SENDER_IDENTITY),
			},
			Actions: []hermes.Action{
				{
					Instructions: "To reset your account password, Please use the code below:",
					InviteCode:   token,
				},
			},
			Outros: []string{
				fmt.Sprintf("The code above is only valid until %s. Please enter the code immediately.", expiredAt),
				"Please do not share the OTP code with anyone to keep your data confidential.",
				"If you didn't make this request, you can safely ignore this mail. No changes will be made to your account information.",
				"This is an automated email, please don't send or reply to this email.",
			},
		},
	}

	htmlString, err := e.generateHTML(&email)

	if err != nil {
		return nil, err
	}

	return htmlString, nil
}

func (e *emailHTMLImpl) generateHTML(email *hermes.Email) (*string, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: config.GetAppConfig().EMAIL_SENDER_IDENTITY,
			// Link: "https://example-tes.com/",
			// Logo: "https://example-tes.com/images/example.png",
			// Custom copyright notice
			Copyright: fmt.Sprintf("Copyright © 2024 %s. All rights reserved.", config.GetAppConfig().EMAIL_SENDER_IDENTITY),
		},
	}

	emailBody, err := h.GenerateHTML(*email)
	if err != nil {
		log.Printf("[GenerateHTML - Helper] err: %s", err.Error())
		return nil, err
	}

	return &emailBody, nil
}
