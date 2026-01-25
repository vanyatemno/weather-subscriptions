package subscriptions

import (
	"errors"
	"testing"
	"time"

	"weather-subscriptions/internal/db/models"

	"gorm.io/gorm"
)

type tokensStateMock struct {
	getTokenFunc                func(token string) (*models.Token, error)
	getTokenByUserIDAndTypeFunc func(userID string, tokenType string) (*models.Token, error)
	saveTokenFunc               func(token *models.Token) error
	removeTokenFunc             func(token *models.Token) error
	savedToken                  *models.Token
	removedToken                *models.Token
}

func (m *tokensStateMock) GetToken(token string) (*models.Token, error) {
	if m.getTokenFunc == nil {
		return nil, nil
	}
	return m.getTokenFunc(token)
}

func (m *tokensStateMock) GetTokenByUserIDAndType(userID string, tokenType string) (*models.Token, error) {
	if m.getTokenByUserIDAndTypeFunc == nil {
		return nil, nil
	}
	return m.getTokenByUserIDAndTypeFunc(userID, tokenType)
}

func (m *tokensStateMock) SaveToken(token *models.Token) error {
	m.savedToken = token
	if m.saveTokenFunc == nil {
		return nil
	}
	return m.saveTokenFunc(token)
}

func (m *tokensStateMock) RemoveToken(token *models.Token) error {
	m.removedToken = token
	if m.removeTokenFunc == nil {
		return nil
	}
	return m.removeTokenFunc(token)
}

//nolint:gocognit
func TestVerifyToken(t *testing.T) {
	now := time.Now()
	validToken := &models.Token{Token: "abc", ExpiryAt: now.Add(time.Hour)}
	expiredToken := &models.Token{Token: "expired", ExpiryAt: now.Add(-time.Hour)}
	deletedToken := &models.Token{
		Token:    "deleted",
		ExpiryAt: now.Add(time.Hour),
		DeletedAt: gorm.DeletedAt{
			Time:  now.Add(-time.Minute),
			Valid: true,
		},
	}

	tests := []struct {
		name      string
		token     string
		getFunc   func(token string) (*models.Token, error)
		wantError bool
	}{
		{
			name:  "success",
			token: "abc",
			getFunc: func(token string) (*models.Token, error) {
				return validToken, nil
			},
		},
		{
			name:  "expired",
			token: "expired",
			getFunc: func(token string) (*models.Token, error) {
				return expiredToken, nil
			},
			wantError: true,
		},
		{
			name:  "deleted",
			token: "deleted",
			getFunc: func(token string) (*models.Token, error) {
				return deletedToken, nil
			},
			wantError: true,
		},
		{
			name:  "state error",
			token: "err",
			getFunc: func(token string) (*models.Token, error) {
				return nil, errors.New("boom")
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenState := &tokensStateMock{getTokenFunc: test.getFunc}
			manager := &SubscriptionManager{tokenState: tokenState}

			result, err := manager.verifyToken(test.token)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if result != nil {
					t.Fatalf("expected nil token, got %v", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatalf("expected token, got nil")
			}
			if result.Token != validToken.Token {
				t.Fatalf("expected token %q, got %q", validToken.Token, result.Token)
			}
		})
	}
}

func TestCheckUserToken(t *testing.T) {
	now := time.Now()
	expiredToken := &models.Token{Token: "expired", ExpiryAt: now.Add(-time.Hour)}
	activeToken := &models.Token{Token: "active", ExpiryAt: now.Add(time.Hour)}

	tests := []struct {
		name          string
		getFunc       func(userID string, tokenType string) (*models.Token, error)
		removeFunc    func(token *models.Token) error
		wantError     bool
		expectRemoved bool
		expectedToken *models.Token
	}{
		{
			name: "record not found",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return nil, gorm.ErrRecordNotFound
			},
		},
		{
			name: "state error",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return nil, errors.New("boom")
			},
			wantError: true,
		},
		{
			name: "expired token returns error",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return expiredToken, nil
			},
			wantError: true,
		},
		{
			name: "active token removed",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return activeToken, nil
			},
			removeFunc: func(token *models.Token) error {
				return nil
			},
			expectRemoved: true,
			expectedToken: activeToken,
		},
		{
			name: "remove error",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return activeToken, nil
			},
			removeFunc: func(token *models.Token) error {
				return errors.New("remove")
			},
			wantError:     true,
			expectRemoved: true,
			expectedToken: activeToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenState := &tokensStateMock{
				getTokenByUserIDAndTypeFunc: test.getFunc,
				removeTokenFunc:             test.removeFunc,
			}
			manager := &SubscriptionManager{tokenState: tokenState}

			err := manager.checkUserToken("user", models.Sub)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if test.expectRemoved {
				if tokenState.removedToken == nil {
					t.Fatalf("expected token to be removed")
				}
				if tokenState.removedToken != test.expectedToken {
					t.Fatalf("expected removed token %v, got %v", test.expectedToken, tokenState.removedToken)
				}
			}
		})
	}
}

//nolint:gocognit
func TestCreateToken(t *testing.T) {
	frequency := models.DAILY

	tests := []struct {
		name        string
		getFunc     func(userID string, tokenType string) (*models.Token, error)
		saveFunc    func(token *models.Token) error
		wantError   bool
		expectSaved bool
	}{
		{
			name: "success",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return nil, gorm.ErrRecordNotFound
			},
			saveFunc: func(token *models.Token) error {
				return nil
			},
			expectSaved: true,
		},
		{
			name: "check user token error",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return nil, errors.New("boom")
			},
			wantError: true,
		},
		{
			name: "save error",
			getFunc: func(userID string, tokenType string) (*models.Token, error) {
				return nil, gorm.ErrRecordNotFound
			},
			saveFunc: func(token *models.Token) error {
				return errors.New("save")
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenState := &tokensStateMock{
				getTokenByUserIDAndTypeFunc: test.getFunc,
				saveTokenFunc:               test.saveFunc,
			}
			manager := &SubscriptionManager{tokenState: tokenState}

			token, err := manager.createToken("user", models.Sub, &frequency)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token == nil {
				t.Fatalf("expected token, got nil")
			}
			if token.Token == "" {
				t.Fatalf("expected generated token to be non-empty")
			}
			if token.Type != models.Sub {
				t.Fatalf("expected token type %q, got %q", models.Sub, token.Type)
			}
			if token.UserID != "user" {
				t.Fatalf("expected user ID 'user', got %q", token.UserID)
			}
			if token.SubscriptionType == nil || *token.SubscriptionType != frequency {
				t.Fatalf("expected subscription type %q", frequency)
			}
			if test.expectSaved {
				if tokenState.savedToken == nil {
					t.Fatalf("expected token to be saved")
				}
				if tokenState.savedToken != token {
					t.Fatalf("expected saved token to match created token")
				}
			}
		})
	}
}
