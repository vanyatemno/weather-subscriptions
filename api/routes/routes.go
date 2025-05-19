package routes

import (
	"github.com/gofiber/fiber/v2"
	"weather-subscriptions/api/handlers"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/mail/mailer_service"
	"weather-subscriptions/internal/state"
)

type Routes struct {
	handler *handlers.RequestHandler
}

func New(cfg *config.Config, state state.Stateful, mailer mailer_service.MailerService) *Routes {
	handler := handlers.New(cfg, state, mailer)
	return &Routes{handler}
}

func (r *Routes) Setup(app *fiber.App) {
	app.Get("/weather", r.handler.WeatherHandler.GetWeather)
	app.Post("/subscribe", r.handler.SubscriptionHandler.HandleSubscribe)
	app.Get("/confirm/:token", r.handler.SubscriptionHandler.HandleConfirmSubscription)
	app.Get("/unsubscribe/:token", r.handler.SubscriptionHandler.HandleUnsubscribe)
}
