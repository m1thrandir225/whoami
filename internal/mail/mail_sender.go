package mail

type MailService interface {
	SendMail(from, to, subject, content string) error
}
