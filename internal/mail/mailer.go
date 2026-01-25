package mail

import (
	"weather-subscriptions/internal/config"

	"gopkg.in/gomail.v2"
)

type MailerService interface {
	Send(message Message) error
}

type Mailer struct {
	cfg *config.Config
}

func NewMailerService(cfg *config.Config) MailerService {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) Send(message Message) error {
	client := gomail.NewDialer(m.cfg.Mailer.SMTP, m.cfg.Mailer.Port, m.cfg.Mailer.From, m.cfg.Mailer.Password)

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.Mailer.From)
	msg.SetHeader("To", message.To[0])
	msg.SetHeader("Subject", message.Subject)
	msg.SetBody("text/html", message.Body)
	if err := client.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
