package subscriptions

import (
	"context"
	"testing"
	"time"
	"weather-subscriptions/internal/db/models"

	// Placeholder for actual mock paths if they exist
	// imocks "weather-subscriptions/internal/integrations/mocks"
	// mmocks "weather-subscriptions/internal/mail/mailer_service/mocks"
	// smocks "weather-subscriptions/internal/state/mocks"
	"weather-subscriptions/internal/integrations"        // Using actual interface for mock struct
	"weather-subscriptions/internal/mail/mailer_service" // Using actual interface for mock struct
	"weather-subscriptions/internal/state"               // Using actual interface for mock struct
	"weather-subscriptions/internal/templates"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockState is a mock for state.Stateful
type MockState struct {
	mock.Mock
	state.Stateful // Embed the interface
}

func (m *MockState) GetCity(name string) (*models.City, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *MockState) SaveCity(city *models.City) error {
	args := m.Called(city)
	return args.Error(0)
}

func (m *MockState) SaveUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockState) GetUser(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockState) RemoveUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockState) SaveToken(token *models.Token) (*models.Token, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockState) GetTokenByValue(value string) (*models.Token, error) {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockState) RemoveToken(token *models.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockState) SaveSubscription(sub *models.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

func (m *MockState) GetSubscription(userID string) (*models.Subscription, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockState) RemoveSubscription(sub *models.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

// MockMapsIntegration is a mock for integrations.MapsIntegration
type MockMapsIntegration struct {
	mock.Mock
	integrations.MapsIntegration // Embed the interface
}

func (m *MockMapsIntegration) GetCity(ctx context.Context, name string) (*models.City, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *MockMapsIntegration) GetWeather(ctx context.Context, lat, lon float64) (*models.Weather, error) {
	args := m.Called(ctx, lat, lon)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Weather), args.Error(1)
}


// MockMailerService is a mock for mailer_service.MailerService
type MockMailerService struct {
	mock.Mock
	mailer_service.MailerService // Embed the interface
}

func (m *MockMailerService) Send(message mailer_service.MailMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMailerService) SendBatch(messages []mailer_service.MailMessage) []error {
	args := m.Called(messages)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]error)
}


// TestSubscriptionManager_SendConfirmationEmail_NewCity_Success tests the happy path where a new city is fetched and saved.
func TestSubscriptionManager_SendConfirmationEmail_NewCity_Success(t *testing.T) {
	mockState := new(MockState)
	mockMaps := new(MockMapsIntegration)
	mockMailer := new(MockMailerService)

	// Ensure all methods of the interfaces are present on the mocks
	// by assigning them to interface variables (optional, for compile-time check)
	var _ state.Stateful = mockState
	var _ integrations.MapsIntegration = mockMaps
	var _ mailer_service.MailerService = mockMailer


	manager := New(mockState, mockMailer, mockMaps).(*SubscriptionManager) 

	request := SubscribeRequest{
		Email:     "test@example.com",
		City:      "NewCity",
		Frequency: "daily",
	}

	ctx := context.Background()
	cityID := uuid.NewString()
	cityModel := &models.City{ID: cityID, Name: "NewCity", Lat: 1.0, Lon: 1.0}
	
	mockState.On("GetCity", request.City).Return(nil, gorm.ErrRecordNotFound).Once()
	mockMaps.On("GetCity", ctx, request.City).Return(cityModel, nil).Once()
	mockState.On("SaveCity", cityModel).Return(nil).Once()
	
	var capturedUserID string
	mockState.On("SaveUser", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		userArg := args.Get(0).(*models.User)
		assert.Equal(t, request.Email, userArg.Email)
		assert.Equal(t, cityModel.ID, userArg.CityID)
		capturedUserID = userArg.ID // Capture the generated user ID
		assert.NotEmpty(t, capturedUserID)
	}).Return(nil).Once()

	expectedSubTokenValue := "subToken123"
	expectedUnsubTokenValue := "unsubToken456"

	mockState.On("SaveToken", mock.MatchedBy(func(token *models.Token) bool {
		return token.UserID == capturedUserID && token.Type == string(models.Sub) && token.SubscriptionType == request.Frequency
	})).Run(func(args mock.Arguments) {
		tokenArg := args.Get(0).(*models.Token)
		tokenArg.Token = expectedSubTokenValue // Simulate token value assignment
	}).Return(&models.Token{Token: expectedSubTokenValue}, nil).Once()
	
	mockState.On("SaveToken", mock.MatchedBy(func(token *models.Token) bool {
		return token.UserID == capturedUserID && token.Type == string(models.Unsub)
	})).Run(func(args mock.Arguments) {
		tokenArg := args.Get(0).(*models.Token)
		tokenArg.Token = expectedUnsubTokenValue // Simulate token value assignment
	}).Return(&models.Token{Token: expectedUnsubTokenValue}, nil).Once()
	
	mockMailer.On("Send", mock.MatchedBy(func(msg mailer_service.MailMessage) bool {
		return msg.To[0] == request.Email && 
		       msg.Subject == "Confirmation code" && 
		       msg.Body == templates.GetVerificationEmailTemplate(expectedSubTokenValue) // Mail should use the sub token
	})).Return(nil).Once()

	err := manager.SendConfirmationEmail(ctx, request)

	assert.NoError(t, err)
	mockState.AssertExpectations(t)
	mockMaps.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}

func TestSubscriptionManager_SendConfirmationEmail_ExistingCity_Success(t *testing.T) {
	mockState := new(MockState)
	mockMaps := new(MockMapsIntegration)
	mockMailer := new(MockMailerService)

	manager := New(mockState, mockMailer, mockMaps).(*SubscriptionManager)

	request := SubscribeRequest{
		Email:     "test2@example.com",
		City:      "ExistingCity",
		Frequency: "weekly",
	}
	ctx := context.Background()
	cityID := uuid.NewString()
	cityModel := &models.City{ID: cityID, Name: "ExistingCity", Lat: 2.0, Lon: 2.0}

	mockState.On("GetCity", request.City).Return(cityModel, nil).Once()
	
	var capturedUserID string
	mockState.On("SaveUser", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		userArg := args.Get(0).(*models.User)
		assert.Equal(t, request.Email, userArg.Email)
		assert.Equal(t, cityModel.ID, userArg.CityID)
		capturedUserID = userArg.ID
	}).Return(nil).Once()

	expectedSubTokenValue := "subToken789"
	mockState.On("SaveToken", mock.MatchedBy(func(token *models.Token) bool {
		return token.UserID == capturedUserID && token.Type == string(models.Sub) && token.SubscriptionType == request.Frequency
	})).Run(func(args mock.Arguments) {
		args.Get(0).(*models.Token).Token = expectedSubTokenValue
	}).Return(&models.Token{Token: expectedSubTokenValue}, nil).Once()

	mockState.On("SaveToken", mock.MatchedBy(func(token *models.Token) bool {
		return token.UserID == capturedUserID && token.Type == string(models.Unsub)
	})).Run(func(args mock.Arguments) {
		args.Get(0).(*models.Token).Token = "unsubTokenABC"
	}).Return(&models.Token{Token: "unsubTokenABC"}, nil).Once()
	
	mockMailer.On("Send", mock.MatchedBy(func(msg mailer_service.MailMessage) bool {
		return msg.To[0] == request.Email && msg.Body == templates.GetVerificationEmailTemplate(expectedSubTokenValue)
	})).Return(nil).Once()

	err := manager.SendConfirmationEmail(ctx, request)

	assert.NoError(t, err)
	mockState.AssertExpectations(t)
	mockMaps.AssertNotCalled(t, "GetCity", mock.Anything, mock.Anything) 
	mockState.AssertNotCalled(t, "SaveCity", mock.Anything)          
	mockMailer.AssertExpectations(t)
}

func TestSubscriptionManager_Subscribe_Success(t *testing.T) {
	mockState := new(MockState)
	mockMaps := new(MockMapsIntegration) 
	mockMailer := new(MockMailerService)

	manager := New(mockState, mockMailer, mockMaps).(*SubscriptionManager)

	tokenValue := "validSubTokenAbc"
	userID := uuid.NewString()
	userToken := &models.Token{
		ID:               uuid.NewString(),
		Token:            tokenValue,
		UserID:           userID,
		Type:             string(models.Sub),
		SubscriptionType: "daily",
		ExpiresAt:        time.Now().Add(1 * time.Hour),
	}

	mockState.On("GetTokenByValue", tokenValue).Return(userToken, nil).Once()
	mockState.On("SaveSubscription", mock.MatchedBy(func(sub *models.Subscription) bool {
		return sub.UserID == userID && sub.Frequency == userToken.SubscriptionType && sub.ID != ""
	})).Return(nil).Once()
	mockState.On("RemoveToken", userToken).Return(nil).Once()

	err := manager.Subscribe(tokenValue)

	assert.NoError(t, err)
	mockState.AssertExpectations(t)
}

func TestSubscriptionManager_Subscribe_InvalidToken_WrongType(t *testing.T) {
	mockState := new(MockState)
	mockMaps := new(MockMapsIntegration)
	mockMailer := new(MockMailerService)

	manager := New(mockState, mockMailer, mockMaps).(*SubscriptionManager)

	tokenValue := "unsubTokenForSubXyz"
	userToken := &models.Token{
		ID:        uuid.NewString(),
		Token:     tokenValue,
		UserID:    uuid.NewString(),
		Type:      string(models.Unsub), 
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockState.On("GetTokenByValue", tokenValue).Return(userToken, nil).Once()

	err := manager.Subscribe(tokenValue)

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
	mockState.AssertCalled(t, "GetTokenByValue", tokenValue)
	mockState.AssertNotCalled(t, "SaveSubscription", mock.Anything)
	// RemoveToken is called in Subscribe *after* successful SaveSubscription.
	// If verifyToken itself removes expired tokens, that would be a different mock on GetTokenByValue.
	// Based on manager.go, RemoveToken in Subscribe is only after SaveSubscription.
	mockState.AssertNotCalled(t, "RemoveToken", mock.AnythingOfType("*models.Token"))
}

func TestSubscriptionManager_Unsubscribe_Success(t *testing.T) {
	mockState := new(MockState)
	mockMaps := new(MockMapsIntegration)
	mockMailer := new(MockMailerService)

	manager := New(mockState, mockMailer, mockMaps).(*SubscriptionManager)

	tokenValue := "validUnsubTokenDef"
	userID := uuid.NewString()
	userToken := &models.Token{
		ID:        uuid.NewString(),
		Token:     tokenValue,
		UserID:    userID,
		Type:      string(models.Unsub),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockState.On("GetTokenByValue", tokenValue).Return(userToken, nil).Once()
	// The RemoveUser call in manager.go passes &models.User{ID: userToken.UserID}
	mockState.On("RemoveUser", mock.MatchedBy(func(u *models.User) bool { return u.ID == userID })).Return(nil).Once()
	mockState.On("RemoveToken", userToken).Return(nil).Once()

	err := manager.Unsubscribe(tokenValue)

	assert.NoError(t, err)
	mockState.AssertExpectations(t)
}
