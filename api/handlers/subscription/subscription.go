package handlers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/mail/mailer_service"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/subscriptions"
)

type SubscriptionHandler struct {
	manager subscriptions.SubManager
}

func NewSubscriptionHandler(
	state state.Stateful,
	mailer mailer_service.MailerService,
	integration integrations.MapsIntegration,
) *SubscriptionHandler {
	manager := subscriptions.New(state, mailer, integration)
	return &SubscriptionHandler{
		manager: manager,
	}
}

// HandleSubscribe handles the POST /subscribe endpoint
func (sh *SubscriptionHandler) HandleSubscribe(c *fiber.Ctx) error {
	var request subscriptions.SubscribeRequest
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(&request)
	if err != nil {
		return err
	}
	request.City = slug.Make(request.City)

	err = sh.manager.InviteUser(c.Context(), request)
	if err != nil && err.Error() == "user already exists" {
		return c.SendStatus(fiber.StatusConflict)
	} else if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "confirmation email sent"})
}

// HandleConfirmSubscription handles the POST /confirm/{token} endpoint
func (sh *SubscriptionHandler) HandleConfirmSubscription(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := sh.manager.Subscribe(token)
	if err != nil && err.Error() == "invalid token" {
		return c.SendStatus(fiber.StatusNotFound)
	} else if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}

// HandleUnsubscribe handles the POST /unsubscribe/{token} endpoint
func (sh *SubscriptionHandler) HandleUnsubscribe(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := sh.manager.Unsubscribe(token)
	if err != nil && err.Error() == "invalid token" {
		return c.SendStatus(fiber.StatusNotFound)
	} else if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}
