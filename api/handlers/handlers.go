package handlers

import (
	subscriptionHandlers "weather-subscriptions/api/handlers/subscription"
	weatherHandlers "weather-subscriptions/api/handlers/weather"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/integrations/google"
	"weather-subscriptions/internal/mail/mailer_service"
	"weather-subscriptions/internal/state"
)

type RequestHandler struct {
	WeatherHandler      *weatherHandlers.WeatherHandler
	SubscriptionHandler *subscriptionHandlers.SubscriptionHandler
}

func New(cfg *config.Config, state state.Stateful, mailer mailer_service.MailerService) *RequestHandler {
	googleInt := google.New(cfg)
	weatherHandler := weatherHandlers.NewWeatherHandler(googleInt, state)
	subscriptionHandler := subscriptionHandlers.NewSubscriptionHandler(state, mailer, googleInt)
	return &RequestHandler{weatherHandler, subscriptionHandler}
}
