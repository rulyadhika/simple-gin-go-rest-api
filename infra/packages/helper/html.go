package helper

import (
	"fmt"
	"log"

	"github.com/matcornic/hermes/v2"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
)

type EmailHTML interface {
	GenerateSendTokenEmailHTML(username string, token string) (*string, error)
}

type emailHTMLImpl struct{}

func NewEmailHTMLImpl() EmailHTML {
	return &emailHTMLImpl{}
}

func (e *emailHTMLImpl) GenerateSendTokenEmailHTML(username string, token string) (*string, error) {
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

func (e *emailHTMLImpl) generateHTML(email *hermes.Email) (*string, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: config.GetAppConfig().EMAIL_SENDER_IDENTITY,
			// Link: "https://example-tes.com/",
			// Logo: "https://example-tes.com/images/example.png",
		},
	}

	emailBody, err := h.GenerateHTML(*email)
	if err != nil {
		log.Printf("[GenerateHTML - Helper] err: %s", err.Error())
		return nil, err
	}

	return &emailBody, nil
}
