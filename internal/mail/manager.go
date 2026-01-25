package mail

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"weather-subscriptions/internal/advises"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/integrations/google"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/templates"
)

const weatherLifetime = 5 * time.Minute

// MailManager interface to send emails to subscribed users
type MailManager interface {
	SendHourly() error
	SendDaily() error
}

type Manager struct {
	ctx                context.Context
	cfg                *config.Config
	subscriptionState  state.SubscriptionsState
	weatherState       state.WeatherState
	citiesState        state.CitiesState
	tokenState         state.TokensState
	mailer             MailerService
	weatherIntegration integrations.MapsIntegration
	advisesService     *advises.AdvisesService
}

// SendHourly sends email with current weather information to users with "hourly" subscription
func (m *Manager) SendHourly() error {
	subscriptions, err := m.subscriptionState.GetSubscriptions(models.HOURLY)
	if err != nil {
		zap.L().Error("failed to get subscriptions", zap.Error(err))
		return err
	}

	err = m.sendMail(subscriptions, models.HOURLY)
	if err != nil {
		zap.L().Error("failed to send mail", zap.Error(err))
		return err
	}

	return nil
}

// SendDaily sends email with current weather information to users with "daily" subscription
func (m *Manager) SendDaily() error {
	subscriptions, err := m.subscriptionState.GetSubscriptions(models.DAILY)
	if err != nil {
		return err
	}

	err = m.sendMail(subscriptions, models.DAILY)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) SendNotification(subscriptionType models.SubscriptionType) error {
	subscriptions, err := m.subscriptionState.GetSubscriptions(subscriptionType)
	if err != nil {
		return err
	}

	err = m.sendMail(subscriptions, subscriptionType)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) sendMail(subscriptions []*models.Subscription, subType models.SubscriptionType) error {
	wg := sync.WaitGroup{}
	for i := range subscriptions {
		weather, err := m.getWeatherForSubscription(subscriptions[i])
		if err != nil {
			zap.L().Error("failed to get weather for subscription", zap.Error(err))
			return err
		}
		unsubToken, err := m.tokenState.GetTokenByUserIDAndType(subscriptions[i].User.ID, models.Unsub)
		if err != nil {
			zap.L().Error("failed to get unsub token", zap.Error(err))
			return err
		}
		advise, err := m.advisesService.GetAdvise(m.ctx, weather)
		if err != nil {
			zap.L().Error("failed to get weather advise for subscription", zap.Error(err))
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = m.mailer.Send(MailMessage{
				To:      []string{subscriptions[i].User.Email},
				Subject: fmt.Sprintf("Your %s weather", strings.ToLower(subType)),
				Body:    templates.GetWeatherEmailBody(weather, m.cfg.FrontendURL, unsubToken.Token, advise),
			})
			if err != nil {
				zap.L().Error("Error sending email", zap.Error(err))
			}
		}()
	}
	wg.Wait()

	return nil
}

func (m *Manager) getWeatherForSubscription(subscription *models.Subscription) (*models.Weather, error) {
	weather, err := m.weatherState.GetWeather(subscription.User.CityID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if weather == nil || weather.Time.Before(time.Now().Add(-weatherLifetime)) {
		city, err := m.citiesState.GetCityByID(subscription.User.CityID)
		if err != nil {
			return nil, err
		}
		weather, err = m.weatherIntegration.GetWeather(m.ctx, city)
		if err != nil {
			return nil, err
		}
	}

	return weather, nil
}

func New(
	ctx context.Context,
	cfg *config.Config,
	subscriptionState state.SubscriptionsState,
	weatherState state.WeatherState,
	citiesState state.CitiesState,
	tokenState state.TokensState,
	mailer MailerService,
	advisesService *advises.AdvisesService,
) MailManager {
	googleInt := google.New(cfg)
	return &Manager{
		cfg:                cfg,
		subscriptionState:  subscriptionState,
		weatherState:       weatherState,
		citiesState:        citiesState,
		tokenState:         tokenState,
		mailer:             mailer,
		ctx:                ctx,
		weatherIntegration: googleInt,
		advisesService:     advisesService,
	}
}
