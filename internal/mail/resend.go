package mail

import (
	gomail "gopkg.in/mail.v2"
)

type ResendMailer struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

func NewResendMailer(smtpHost string, smtpPort int, smtpUsername, smtpPassword string) *ResendMailer {
	return &ResendMailer{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	}
}

func (mail *ResendMailer) SendMail(from string, to string, subject, content string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", content)

	dialer := gomail.NewDialer(
		mail.SMTPHost,
		mail.SMTPPort,
		mail.SMTPUsername,
		mail.SMTPPassword,
	)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
