package mailer_service

// SMTPConfig holds the configuration required to connect and authenticate with an SMTP server.
type SMTPConfig struct {
	Host     string // SMTP server hostname or IP address (e.g., "smtp.example.com").
	Port     int    // SMTP server port (e.g., 587 for TLS, 465 for SMTPS, 25 for unencrypted).
	Username string // Username for SMTP authentication.
	Password string // Password for SMTP authentication.
	From     string // Default sender email address (e.g., "noreply@example.com").
}

// MailMessage represents an email message.
// It contains all necessary fields for constructing and sending an email.
type MailMessage struct {
	To      []string // Primary recipients of the email.
	Subject string   // Subject line of the email.
	Body    string   // Body content of the email (assumed to be plain text).
}
