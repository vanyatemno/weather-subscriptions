package mailer_service

// SMTPConfig holds the configuration required to connect and authenticate with an SMTP server.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// MailMessage represents an email message.
// It contains all necessary fields for constructing and sending an email.
type MailMessage struct {
	To      []string
	Subject string
	Body    string
}
