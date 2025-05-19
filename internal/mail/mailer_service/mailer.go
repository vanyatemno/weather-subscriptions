package mailer_service

import (
	"gopkg.in/gomail.v2"
	"weather-subscriptions/internal/config"
)

type MailerService interface {
	Send(message MailMessage) error
}

type Mailer struct {
	cfg *config.Config
}

func New(cfg *config.Config) MailerService {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) Send(message MailMessage) error {
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
