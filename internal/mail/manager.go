package mail

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/integrations"
	mail "weather-subscriptions/internal/mail/mailer_service"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/templates"
)

const weatherLifetime = 5 * time.Minute

type MailManager interface {
	SendHourly() error
	SendDaily() error
}

type Manager struct {
	cfg                *config.Config
	state              state.Stateful
	mailer             mail.MailerService
	weatherIntegration integrations.MapsIntegration
	ctx                context.Context
}

func (m *Manager) SendHourly() error {
	subscriptions, err := m.state.GetSubscriptions(models.HOURLY)
	if err != nil {
		return err
	}

	err = m.sendMail(subscriptions, models.HOURLY)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) SendDaily() error {
	subscriptions, err := m.state.GetSubscriptions(models.DAILY)
	if err != nil {
		return err
	}

	err = m.sendMail(subscriptions, models.DAILY)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) sendMail(subscriptions []*models.Subscription, subType models.SubscriptionType) error {
	wg := sync.WaitGroup{}
	for i := range subscriptions {
		weather, err := m.state.GetWeather(subscriptions[i].User.CityID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if weather == nil || weather.Time.Before(time.Now().Add(-weatherLifetime)) {
			city, err := m.state.GetCityByID(subscriptions[i].User.CityID)
			if err != nil {
				return err
			}
			weather, err = m.weatherIntegration.GetWeather(m.ctx, city)
			if err != nil {
				return err
			}
		}
		unsubToken, err := m.state.GetToken(subscriptions[i].User.ID)

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.mailer.Send(mail.MailMessage{
				To:      []string{subscriptions[i].User.Email},
				Subject: fmt.Sprintf("Your %s weathre", strings.ToLower(string(subType))),
				Body:    templates.GetWeatherEmailBody(weather, unsubToken.Token),
			})
		}()

	}

	return nil
}

func New(ctx context.Context, cfg *config.Config, state state.Stateful, mailer mail.MailerService) MailManager {
	return &Manager{
		cfg:    cfg,
		state:  state,
		mailer: mailer,
		ctx:    ctx,
	}
}
