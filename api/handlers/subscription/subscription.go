package handlers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"

	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/mail"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/subscriptions"
)

type SubscriptionHandler struct {
	manager subscriptions.SubManager
}

type SubscribeResponse struct {
	Message string `json:"message" example:"confirmation email sent"`
}

func NewSubscriptionHandler(
	cfg *config.Config,
	state *state.State,
	mailer mail.MailerService,
	integration integrations.MapsIntegration,
) *SubscriptionHandler {
	manager := subscriptions.New(
		cfg,
		state,
		state,
		state,
		state,
		mailer,
		integration,
	)
	return &SubscriptionHandler{
		manager: manager,
	}
}

// HandleSubscribe handles the POST /subscribe endpoint
// @Summary Subscribe to weather updates
// @Description Subscribe an email to receive weather updates for a specific city and frequency.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body subscriptions.SubscribeRequest true "Subscription request"
// @Success 200 {object} SubscribeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /subscribe [post]
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
// @Summary Confirm email subscription
// @Description Confirms a subscription using the token sent in the confirmation email.
// @Tags subscriptions
// @Param token path string true "Confirmation token"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /confirm/{token} [get]
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
// @Summary Unsubscribe from weather updates
// @Description Unsubscribes an email from weather updates using the token sent in emails.
// @Tags subscriptions
// @Param token path string true "Unsubscribe token"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /unsubscribe/{token} [get]
func (sh *SubscriptionHandler) HandleUnsubscribe(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := sh.manager.Unsubscribe(token)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendStatus(fiber.StatusOK)
}
