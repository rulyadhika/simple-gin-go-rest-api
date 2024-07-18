package helper

import (
	"fmt"
	"log"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
)

func SendActivationTokenEmail(data dto.SendActivationTokenEmailRequest) error {
	htmlEmailBody, err := NewEmailHTMLImpl().GenerateSendActivationTokenEmailHTML(data.Username, data.Token)

	if err != nil {
		return err
	}

	if err := sendMail(data.Subject, []string{data.ToEmailAddress}, []byte(*htmlEmailBody)); err != nil {
		return err
	}

	return nil
}

func SendPasswordResetTokenEmail(data dto.SendResetPasswordTokenEmailRequest) error {
	htmlEmailBody, err := NewEmailHTMLImpl().GenerateSendPasswordResetTokenEmailHTML(data.Username, data.Token, data.ExpiredAt)

	if err != nil {
		return err
	}

	if err := sendMail(data.Subject, []string{data.ToEmailAddress}, []byte(*htmlEmailBody)); err != nil {
		return err
	}

	return nil
}

func sendMail(subject string, to []string, htmlEmailBody []byte) error {
	appConfig := config.GetAppConfig()

	em := &email.Email{
		To:      to,
		From:    fmt.Sprintf("%s <%s>", appConfig.EMAIL_SENDER_IDENTITY, appConfig.EMAIL_SENDER_AND_SMTP_USER),
		Subject: fmt.Sprintf("[%s] %s", appConfig.EMAIL_SENDER_IDENTITY, subject),
		HTML:    htmlEmailBody,
		Headers: textproto.MIMEHeader{},
	}

	if err := em.Send(fmt.Sprintf("%s:%v", appConfig.EMAIL_SMTP_SERVER, appConfig.EMAIL_SMTP_PORT), smtp.PlainAuth("", appConfig.EMAIL_SENDER_AND_SMTP_USER, appConfig.EMAIL_SMTP_USER_PASSWORD, appConfig.EMAIL_SMTP_SERVER)); err != nil {
		log.Printf("[SendMail - Helper] err: %s", err.Error())
		return err
	}

	return nil
}
