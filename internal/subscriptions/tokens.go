package subscriptions

import (
	"errors"
	"time"

	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/util"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *SubscriptionManager) verifyToken(token string) (*models.Token, error) {
	foundToken, err := s.tokenState.GetToken(token)
	if err != nil {
		return nil, err
	}
	if foundToken.ExpiryAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	if foundToken.DeletedAt.Valid && foundToken.DeletedAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}
	return foundToken, nil
}

func (s *SubscriptionManager) createToken(
	userID string,
	tokenType models.TokenType,
	frequency *models.SubscriptionType,
) (*models.Token, error) {
	code, err := util.GenerateCode(emailValidationCodeLength)
	if err != nil {
		zap.L().Error("failed to generate code", zap.Error(err))
		return nil, errors.New("failed to generate code")
	}

	err = s.checkUserToken(userID, tokenType)
	if err != nil {
		zap.L().Error("failed to check user token existence", zap.Error(err))
		return nil, err
	}

	token := &models.Token{
		Token:            code,
		Type:             tokenType,
		SubscriptionType: frequency,
		UserID:           userID,
	}
	token.SetTokenExpiry()

	err = s.tokenState.SaveToken(token)
	if err != nil {
		zap.L().Error("failed to save token", zap.Error(err))
		return nil, errors.New("failed to save token")
	}

	return token, nil
}

// checkUserToken - checks if token exists for specified user.
// Removes token if it is expired. Returns an error if active token is present.
func (s *SubscriptionManager) checkUserToken(userID string, tokenType models.TokenType) error {
	token, err := s.tokenState.GetTokenByUserIDAndType(userID, tokenType)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if token.ExpiryAt.Before(time.Now()) {
		return errors.New("token already exists")
	}
	err = s.tokenState.RemoveToken(token)
	if err != nil {
		return err
	}

	return nil
}
