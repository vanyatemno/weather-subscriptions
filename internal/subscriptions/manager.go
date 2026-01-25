package subscriptions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"

	"go.uber.org/zap"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/mail"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/templates"
)

const emailValidationCodeLength = 6

type SubManager interface {
	InviteUser(ctx context.Context, request SubscribeRequest) error
	Subscribe(token string) error
	Unsubscribe(token string) error
}

type SubscribeRequest struct {
	Email     string `validate:"required,email" json:"email" form:"email"`
	City      string `validate:"required" json:"city" form:"city"`
	Frequency string `validate:"required" json:"frequency" form:"frequency"`
}

type SubscriptionManager struct {
	cfg                *config.Config
	citiesState        state.CitiesState
	usersState         state.UsersState
	tokenState         state.TokensState
	subscriptionsState state.SubscriptionsState
	mapsIntegration    integrations.MapsIntegration
	mailer             mail.MailerService
}

func New(
	config *config.Config,
	citiesState state.CitiesState,
	usersState state.UsersState,
	tokenState state.TokensState,
	subscriptionsState state.SubscriptionsState,
	mailerService mail.MailerService,
	integration integrations.MapsIntegration,
) SubManager {
	return &SubscriptionManager{
		cfg:                config,
		citiesState:        citiesState,
		usersState:         usersState,
		tokenState:         tokenState,
		subscriptionsState: subscriptionsState,
		mailer:             mailerService,
		mapsIntegration:    integration,
	}
}

// InviteUser accepts user request for subscription, finds or creates city, creates user record,
// creates confirmation token and sends it to user email
func (s *SubscriptionManager) InviteUser(ctx context.Context, request SubscribeRequest) error {
	city, err := s.citiesState.GetCity(request.City)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		city, err = s.mapsIntegration.GetCity(ctx, slug.Make(request.City))
		if err != nil {
			zap.L().Error("error getting city", zap.Error(err))
			return err
		}
		err = s.citiesState.SaveCity(city)
		if err != nil {
			zap.L().Error("error saving city", zap.Error(err))
			return err
		}
	} else if err != nil {
		zap.L().Error("error getting city", zap.Error(err))
		return err
	}

	user, err := s.usersState.GetUserByEmail(request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if user != nil {
		return errors.New("user already exists")
	}
	user = &models.User{
		ID:     uuid.Must(uuid.NewV7()).String(),
		Email:  request.Email,
		CityID: city.ID,
		City:   *city,
	}
	err = s.usersState.SaveUser(user)
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

	err = s.mailer.Send(mail.MailMessage{
		To:      []string{user.Email},
		Subject: "Confirmation code",
		Body:    templates.GetVerificationEmailTemplate(s.cfg.FrontendURL, token.Token),
	})
	if err != nil {
		zap.L().Error("error sending confirmation email", zap.Error(err))
		return err
	}

	return nil
}

// Subscribe checks if sub token exists and creates subscription for the user
func (s *SubscriptionManager) Subscribe(token string) error {
	userToken, err := s.verifyToken(token)
	if err != nil {
		return errors.New("invalid token")
	}
	if userToken.Type != models.Sub {
		zap.L().Error("found a token of invalid type", zap.String("token", token))
		return errors.New("invalid token")
	}

	subscription := &models.Subscription{
		ID:        uuid.Must(uuid.NewV7()).String(),
		Frequency: *userToken.SubscriptionType,
		UserID:    userToken.UserID,
	}
	err = s.subscriptionsState.SaveSubscription(subscription)
	if err != nil {
		zap.L().Error("error saving subscription", zap.Error(err))
		return err
	}
	err = s.tokenState.RemoveToken(userToken)
	if err != nil {
		zap.L().Error("error removing token", zap.Error(err))
		return err
	}

	return nil
}

// Unsubscribe checks if unsub token exists and deletes user, and all related records
func (s *SubscriptionManager) Unsubscribe(token string) error {
	userToken, err := s.verifyToken(token)
	if err != nil {
		return errors.New("invalid token")
	}
	if userToken.Type != models.Unsub {
		return errors.New("invalid token")
	}

	err = s.usersState.RemoveUser(&models.User{ID: userToken.UserID})
	if err != nil {
		zap.L().Error("error removing user", zap.Error(err))
		return err
	}
	err = s.tokenState.RemoveToken(userToken)
	if err != nil {
		zap.L().Error("error removing token", zap.Error(err))
		return err
	}

	return nil
}
