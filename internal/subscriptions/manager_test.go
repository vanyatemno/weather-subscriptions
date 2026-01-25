package subscriptions

import (
	"context"
	"errors"
	"testing"
	"time"

	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/mail"

	"gorm.io/gorm"
)

type citiesStateMock struct {
	getCityFunc  func(name string) (*models.City, error)
	saveCityFunc func(city *models.City) error
	savedCity    *models.City
}

func (m *citiesStateMock) GetCity(name string) (*models.City, error) {
	if m.getCityFunc == nil {
		return nil, nil
	}
	return m.getCityFunc(name)
}

func (m *citiesStateMock) GetCityByID(_ string) (*models.City, error) {
	return nil, nil
}

func (m *citiesStateMock) SaveCity(city *models.City) error {
	m.savedCity = city
	if m.saveCityFunc == nil {
		return nil
	}
	return m.saveCityFunc(city)
}

type usersStateMock struct {
	getUserByEmailFunc func(email string) (*models.User, error)
	saveUserFunc       func(user *models.User) error
	removeUserFunc     func(user *models.User) error
	savedUser          *models.User
	removedUser        *models.User
}

func (m *usersStateMock) GetUser(_ string) (*models.User, error) {
	return nil, nil
}

func (m *usersStateMock) GetUserByEmail(email string) (*models.User, error) {
	if m.getUserByEmailFunc == nil {
		return nil, nil
	}
	return m.getUserByEmailFunc(email)
}

func (m *usersStateMock) SaveUser(user *models.User) error {
	m.savedUser = user
	if m.saveUserFunc == nil {
		return nil
	}
	return m.saveUserFunc(user)
}

func (m *usersStateMock) RemoveUser(user *models.User) error {
	m.removedUser = user
	if m.removeUserFunc == nil {
		return nil
	}
	return m.removeUserFunc(user)
}

type subscriptionsStateMock struct {
	saveFunc func(subscription *models.Subscription) error
	savedSub *models.Subscription
}

func (m *subscriptionsStateMock) GetSubscription(_ string) (*models.Subscription, error) {
	return nil, nil
}

func (m *subscriptionsStateMock) GetSubscriptions(
	_ models.SubscriptionType,
) ([]*models.Subscription, error) {
	return nil, nil
}

func (m *subscriptionsStateMock) SaveSubscription(subscription *models.Subscription) error {
	m.savedSub = subscription
	if m.saveFunc == nil {
		return nil
	}
	return m.saveFunc(subscription)
}

func (m *subscriptionsStateMock) RemoveSubscription(_ *models.Subscription) error {
	return nil
}

type mapsIntegrationMock struct {
	getCityFunc func(ctx context.Context, cityName string) (*models.City, error)
}

func (m *mapsIntegrationMock) GetWeather(ctx context.Context, city *models.City) (*models.Weather, error) {
	return nil, nil
}

func (m *mapsIntegrationMock) GetCity(ctx context.Context, cityName string) (*models.City, error) {
	if m.getCityFunc == nil {
		return nil, nil
	}
	return m.getCityFunc(ctx, cityName)
}

type mailerMock struct {
	sendFunc func(message mail.Message) error
	sent     []mail.Message
}

func (m *mailerMock) Send(message mail.Message) error {
	m.sent = append(m.sent, message)
	if m.sendFunc == nil {
		return nil
	}
	return m.sendFunc(message)
}

//nolint:gocognit
func TestInviteUser(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{FrontendURL: "http://localhost"}
	city := &models.City{ID: "city-1", Name: "Warsaw"}
	request := SubscribeRequest{Email: "user@example.com", City: "Warsaw", Frequency: models.DAILY}

	tests := []struct {
		name            string
		citiesState     *citiesStateMock
		usersState      *usersStateMock
		tokensState     *tokensStateMock
		mapsIntegration *mapsIntegrationMock
		mailer          *mailerMock
		wantError       bool
		assert          func(t *testing.T, tokensSaved int, usersState *usersStateMock, citiesState *citiesStateMock, mailer *mailerMock)
	}{
		{
			name: "success with existing city",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			tokensState: &tokensStateMock{
				getTokenByUserIDAndTypeFunc: func(userID string, tokenType string) (*models.Token, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			assert: func(t *testing.T, tokensSaved int, usersState *usersStateMock, citiesState *citiesStateMock, mailer *mailerMock) {
				if usersState.savedUser == nil {
					t.Fatalf("expected user to be saved")
				}
				if usersState.savedUser.Email != request.Email {
					t.Fatalf("expected email %q, got %q", request.Email, usersState.savedUser.Email)
				}
				if citiesState.savedCity != nil {
					t.Fatalf("did not expect city to be saved when already found")
				}
				if tokensSaved != 2 {
					t.Fatalf("expected 2 tokens saved, got %d", tokensSaved)
				}
				if len(mailer.sent) != 1 {
					t.Fatalf("expected 1 email sent, got %d", len(mailer.sent))
				}
			},
		},
		{
			name: "city not found uses integration",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			tokensState: &tokensStateMock{
				getTokenByUserIDAndTypeFunc: func(userID string, tokenType string) (*models.Token, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			mapsIntegration: &mapsIntegrationMock{
				getCityFunc: func(ctx context.Context, cityName string) (*models.City, error) {
					return city, nil
				},
			},
			mailer: &mailerMock{},
			assert: func(t *testing.T, tokensSaved int, usersState *usersStateMock, citiesState *citiesStateMock, mailer *mailerMock) {
				if citiesState.savedCity == nil {
					t.Fatalf("expected city to be saved")
				}
			},
		},
		{
			name: "get city error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return nil, errors.New("city")
				},
			},
			usersState:      &usersStateMock{},
			tokensState:     &tokensStateMock{},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			wantError:       true,
		},
		{
			name: "integration get city error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			usersState:  &usersStateMock{},
			tokensState: &tokensStateMock{},
			mapsIntegration: &mapsIntegrationMock{
				getCityFunc: func(ctx context.Context, cityName string) (*models.City, error) {
					return nil, errors.New("integration")
				},
			},
			mailer:    &mailerMock{},
			wantError: true,
		},
		{
			name: "save city error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return nil, gorm.ErrRecordNotFound
				},
				saveCityFunc: func(city *models.City) error {
					return errors.New("save city")
				},
			},
			usersState:  &usersStateMock{},
			tokensState: &tokensStateMock{},
			mapsIntegration: &mapsIntegrationMock{
				getCityFunc: func(ctx context.Context, cityName string) (*models.City, error) {
					return city, nil
				},
			},
			mailer:    &mailerMock{},
			wantError: true,
		},
		{
			name: "get user error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, errors.New("user")
				},
			},
			tokensState:     &tokensStateMock{},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			wantError:       true,
		},
		{
			name: "user already exists",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return &models.User{Email: email}, nil
				},
			},
			tokensState:     &tokensStateMock{},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			wantError:       true,
		},
		{
			name: "save user error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
				saveUserFunc: func(user *models.User) error {
					return errors.New("save user")
				},
			},
			tokensState:     &tokensStateMock{},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			wantError:       true,
		},
		{
			name: "create sub token error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			tokensState: &tokensStateMock{
				getTokenByUserIDAndTypeFunc: func(userID string, tokenType string) (*models.Token, error) {
					return nil, errors.New("token")
				},
			},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			wantError:       true,
		},
		{
			name: "create unsub token error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			tokensState: &tokensStateMock{
				getTokenByUserIDAndTypeFunc: func(userID string, tokenType string) (*models.Token, error) {
					return nil, gorm.ErrRecordNotFound
				},
				saveTokenFunc: func(token *models.Token) error {
					return errors.New("save token")
				},
			},
			mapsIntegration: &mapsIntegrationMock{},
			mailer:          &mailerMock{},
			wantError:       true,
		},
		{
			name: "mailer error",
			citiesState: &citiesStateMock{
				getCityFunc: func(name string) (*models.City, error) {
					return city, nil
				},
			},
			usersState: &usersStateMock{
				getUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			tokensState: &tokensStateMock{
				getTokenByUserIDAndTypeFunc: func(userID string, tokenType string) (*models.Token, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			mapsIntegration: &mapsIntegrationMock{},
			mailer: &mailerMock{
				sendFunc: func(message mail.Message) error {
					return errors.New("send")
				},
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			savedTokens := 0
			if test.tokensState != nil {
				originalSave := test.tokensState.saveTokenFunc
				test.tokensState.saveTokenFunc = func(token *models.Token) error {
					savedTokens++
					if originalSave != nil {
						return originalSave(token)
					}
					return nil
				}
			}

			manager := &SubscriptionManager{
				cfg:                cfg,
				citiesState:        test.citiesState,
				usersState:         test.usersState,
				tokenState:         test.tokensState,
				subscriptionsState: &subscriptionsStateMock{},
				mapsIntegration:    test.mapsIntegration,
				mailer:             test.mailer,
			}

			err := manager.InviteUser(ctx, request)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.assert != nil {
				test.assert(t, savedTokens, test.usersState, test.citiesState, test.mailer)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	frequency := models.DAILY
	now := time.Now().Add(time.Hour)

	tests := []struct {
		name            string
		getTokenFunc    func(token string) (*models.Token, error)
		saveSubFunc     func(subscription *models.Subscription) error
		removeTokenFunc func(token *models.Token) error
		wantError       bool
		assert          func(t *testing.T, subscriptionsState *subscriptionsStateMock, tokenState *tokensStateMock)
	}{
		{
			name: "invalid token",
			getTokenFunc: func(token string) (*models.Token, error) {
				return nil, errors.New("bad token")
			},
			wantError: true,
		},
		{
			name: "invalid token type",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{Token: token, Type: models.Unsub, ExpiryAt: now}, nil
			},
			wantError: true,
		},
		{
			name: "save subscription error",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{
					Token:            token,
					Type:             models.Sub,
					UserID:           "user",
					SubscriptionType: &frequency,
					ExpiryAt:         now,
				}, nil
			},
			saveSubFunc: func(subscription *models.Subscription) error {
				return errors.New("save")
			},
			wantError: true,
		},
		{
			name: "remove token error",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{
					Token:            token,
					Type:             models.Sub,
					UserID:           "user",
					SubscriptionType: &frequency,
					ExpiryAt:         now,
				}, nil
			},
			removeTokenFunc: func(token *models.Token) error {
				return errors.New("remove")
			},
			wantError: true,
		},
		{
			name: "success",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{
					Token:            token,
					Type:             models.Sub,
					UserID:           "user",
					SubscriptionType: &frequency,
					ExpiryAt:         now,
				}, nil
			},
			assert: func(t *testing.T, subscriptionsState *subscriptionsStateMock, tokenState *tokensStateMock) {
				if subscriptionsState.savedSub == nil {
					t.Fatalf("expected subscription to be saved")
				}
				if tokenState.removedToken == nil {
					t.Fatalf("expected token to be removed")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			subState := &subscriptionsStateMock{saveFunc: test.saveSubFunc}
			tokenState := &tokensStateMock{
				getTokenFunc:    test.getTokenFunc,
				removeTokenFunc: test.removeTokenFunc,
			}
			manager := &SubscriptionManager{tokenState: tokenState, subscriptionsState: subState}

			err := manager.Subscribe("token")
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.assert != nil {
				test.assert(t, subState, tokenState)
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	now := time.Now().Add(time.Hour)

	tests := []struct {
		name            string
		getTokenFunc    func(token string) (*models.Token, error)
		removeUserFunc  func(user *models.User) error
		removeTokenFunc func(token *models.Token) error
		wantError       bool
		assert          func(t *testing.T, usersState *usersStateMock, tokenState *tokensStateMock)
	}{
		{
			name: "invalid token",
			getTokenFunc: func(token string) (*models.Token, error) {
				return nil, errors.New("bad")
			},
			wantError: true,
		},
		{
			name: "invalid token type",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{Token: token, Type: models.Sub, UserID: "user", ExpiryAt: now}, nil
			},
			wantError: true,
		},
		{
			name: "remove user error",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{Token: token, Type: models.Unsub, UserID: "user", ExpiryAt: now}, nil
			},
			removeUserFunc: func(user *models.User) error {
				return errors.New("remove user")
			},
			wantError: true,
		},
		{
			name: "remove token error",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{Token: token, Type: models.Unsub, UserID: "user", ExpiryAt: now}, nil
			},
			removeTokenFunc: func(token *models.Token) error {
				return errors.New("remove token")
			},
			wantError: true,
		},
		{
			name: "success",
			getTokenFunc: func(token string) (*models.Token, error) {
				return &models.Token{Token: token, Type: models.Unsub, UserID: "user", ExpiryAt: now}, nil
			},
			assert: func(t *testing.T, usersState *usersStateMock, tokenState *tokensStateMock) {
				if usersState.removedUser == nil {
					t.Fatalf("expected user to be removed")
				}
				if tokenState.removedToken == nil {
					t.Fatalf("expected token to be removed")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usersState := &usersStateMock{removeUserFunc: test.removeUserFunc}
			tokenState := &tokensStateMock{
				getTokenFunc:    test.getTokenFunc,
				removeTokenFunc: test.removeTokenFunc,
			}
			manager := &SubscriptionManager{usersState: usersState, tokenState: tokenState}

			err := manager.Unsubscribe("token")
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.assert != nil {
				test.assert(t, usersState, tokenState)
			}
		})
	}
}
