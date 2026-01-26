package mock

import (
	"weather-subscriptions/internal/mail"
)

type MailerServiceMock struct{}

func NewMailerServiceMock() *MailerServiceMock {
	return &MailerServiceMock{}
}

func (m *MailerServiceMock) Send(_ mail.Message) error {
	return nil
}
