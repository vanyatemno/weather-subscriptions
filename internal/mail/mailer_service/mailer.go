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

//func (m *Mailer) Send(message MailMessage) error {
//	if m.cfg.Mailer.Host == "" || m.cfg.Mailer.Port == 0 || m.cfg.Mailer.From == "" {
//		return fmt.Errorf("mailer configuration is incomplete: host, port, or from address is missing")
//	}
//
//	addr := fmt.Sprintf("%s:%d", m.cfg.Mailer.Host, m.cfg.Mailer.Port)
//
//	// Construct the email message
//	// "From", "To", "Subject" headers, followed by a blank line, then the body.
//	var msgBuilder strings.Builder
//	msgBuilder.WriteString(fmt.Sprintf("From: %s\r\n", m.cfg.Mailer.From))
//	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ",")))
//	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
//	// Assuming HTML body, set Content-Type. For plain text, this can be "text/plain".
//	msgBuilder.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n")
//	msgBuilder.WriteString("\r\n") // Blank line between headers and body
//	msgBuilder.WriteString(message.Body)
//
//	msg := []byte(msgBuilder.String())
//
//	var auth smtp.Auth
//	// Use PlainAuth if username and password are provided
//	if m.cfg.Mailer.Username != "" && m.cfg.Mailer.Password != "" {
//		auth = smtp.PlainAuth("", m.cfg.Mailer.Username, m.cfg.Mailer.Password, m.cfg.Mailer.Host)
//	}
//
//	// Send the email
//	// Note: smtp.SendMail will try to use STARTTLS if the server supports it.
//	// For explicit SSL/TLS, a different approach with `tls.Dial` and `smtp.NewClient` would be needed.
//	// This implementation assumes standard SMTP, potentially with opportunistic STARTTLS.
//	fmt.Println("sending email")
//	err := smtp.SendMail(addr, auth, m.cfg.Mailer.From, message.To, msg)
//	if err != nil {
//		return fmt.Errorf("failed to send email: %w", err)
//	}
//
//	return nil
//}

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
