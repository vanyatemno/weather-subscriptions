package subscriptions

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/integrations"
	mailer2 "weather-subscriptions/internal/mail/mailer_service"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/templates"
)

const emailValidationCodeLength = 6

type SubManager interface {
	SendConfirmationEmail(ctx context.Context, request SubscribeRequest) error
	Subscribe(token string) error
	Unsubscribe(token string) error
}

type SubscribeRequest struct {
	Email     string `validate:"required,email" json:"email" form:"email"`
	City      string `validate:"required" json:"city" form:"city"`
	Frequency string `validate:"required" json:"frequency" form:"frequency"`
}

type SubscriptionManager struct {
	state           state.Stateful
	mapsIntegration integrations.MapsIntegration
	mailer          mailer2.MailerService
}

func New(state state.Stateful, mailer mailer2.MailerService, integration integrations.MapsIntegration) SubManager {
	return &SubscriptionManager{
		state:           state,
		mailer:          mailer,
		mapsIntegration: integration,
	}
}

// SendConfirmationEmail accepts user request for subscription, finds or creates city, creates user record,
// creates confirmation token and sends it to user email
func (s *SubscriptionManager) SendConfirmationEmail(ctx context.Context, request SubscribeRequest) error {
	city, err := s.state.GetCity(request.City)
	if err != nil && errors.Is(gorm.ErrRecordNotFound, err) { // Make sure models is imported if not already
		city, err = s.mapsIntegration.GetCity(ctx, request.City)
		if err != nil {
			zap.L().Error("error getting city", zap.Error(err))
			return err
		}
		err = s.state.SaveCity(city)
		if err != nil {
			zap.L().Error("error saving city", zap.Error(err))
			return err
		}
	} else if err != nil {
		zap.L().Error("error getting city", zap.Error(err))
		return err
	}

	user := &models.User{
		ID:     uuid.Must(uuid.NewV7()).String(),
		Email:  request.Email,
		CityID: city.ID,
		City:   *city,
	}
	err = s.state.SaveUser(user)
	if err != nil {
		zap.L().Error("error saving user", zap.Error(err))
		return err
	}

	// create confirmation code
	token, err := s.createToken(user.ID, models.Sub, &request.Frequency)
	if err != nil {
		zap.L().Error("error creating sub token", zap.Error(err))
		return err
	}
	// create code to unsubscribe
	_, err = s.createToken(user.ID, models.Unsub, nil)
	if err != nil {
		zap.L().Error("error creating unsub token", zap.Error(err))
		return err
	}

	err = s.mailer.Send(mailer2.MailMessage{
		To:      []string{user.Email},
		Subject: "Confirmation code",
		Body:    templates.GetVerificationEmailTemplate(token.Token),
	})
	if err != nil {
		zap.L().Error("error sending confirmation email", zap.Error(err))
		return err
	}

	return nil
}

func (s *SubscriptionManager) Subscribe(token string) error {
	userToken, err := s.verifyToken(token)
	if err != nil {
		return errors.New("invalid token")
	}
	if userToken.Type != string(models.Sub) {
		return errors.New("invalid token")
	}

	subscription := &models.Subscription{
		ID:        uuid.Must(uuid.NewV7()).String(),
		Frequency: userToken.SubscriptionType,
		UserID:    userToken.UserID,
	}
	err = s.state.SaveSubscription(subscription)
	if err != nil {
		return err
	}
	err = s.state.RemoveToken(userToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionManager) Unsubscribe(token string) error {
	userToken, err := s.verifyToken(token)
	if err != nil {
		return errors.New("invalid token")
	}
	if userToken.Type != string(models.Unsub) {
		return errors.New("invalid token")
	}

	err = s.state.RemoveUser(&models.User{ID: userToken.UserID})
	if err != nil {
		return err
	}
	err = s.state.RemoveToken(userToken)
	if err != nil {
		return err
	}

	return nil
}
