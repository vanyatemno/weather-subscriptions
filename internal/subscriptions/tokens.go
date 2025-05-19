package subscriptions

import (
	"crypto/rand"
	"errors"
	"gorm.io/gorm"
	"time"
	"weather-subscriptions/internal/db/models"
)

const (
	subTokenDuration   = 24 * time.Hour
	unsubTokenDuration = time.Hour * 24 * 28 * 13 * 100
)

func (s *SubscriptionManager) verifyToken(token string) (*models.Token, error) {
	foundToken, err := s.state.GetToken(token)
	if err != nil {
		return nil, err
	}
	if foundToken.ExpiryAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return foundToken, nil
}

func (s *SubscriptionManager) createToken(userID string, tokenType models.TokenType, frequency *string) (*models.Token, error) {
	code, err := generateCode()
	if err != nil {
		return nil, errors.New("failed to generate code")
	}

	var token *models.Token
	if tokenType == models.Sub {
		foundToken, err := s.state.GetSubToken(userID)
		if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, err
		}
		if foundToken != nil && foundToken.ExpiryAt.Add(subTokenDuration).Before(time.Now()) {
			return nil, errors.New("token already exists")
		} else if foundToken != nil {
			err = s.state.RemoveToken(foundToken)
			if err != nil {
				return nil, err
			}
		}

		token = &models.Token{
			Token:            code,
			Type:             string(tokenType),
			SubscriptionType: *frequency,
			ExpiryAt:         time.Now().Add(subTokenDuration),
			UserID:           userID,
		}
	} else {
		foundToken, err := s.state.GetUnsubToken(userID)
		if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, err
		}
		if foundToken != nil && foundToken.ExpiryAt.Add(subTokenDuration).Before(time.Now()) {
			return nil, errors.New("token already exists")
		} else if foundToken != nil {
			err = s.state.RemoveToken(foundToken)
			if err != nil {
				return nil, err
			}
		}

		token = &models.Token{
			Token:    code,
			Type:     string(tokenType),
			ExpiryAt: time.Now().Add(unsubTokenDuration),
			UserID:   userID,
		}
	}

	err = s.state.SaveToken(token)
	if err != nil {
		return nil, errors.New("failed to save token")
	}

	return token, nil
}

func generateCode() (string, error) {
	codes := make([]byte, emailValidationCodeLength)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}

	for i := 0; i < emailValidationCodeLength; i++ {
		codes[i] = 48 + (codes[i] % 10)
	}

	return string(codes), nil
}
