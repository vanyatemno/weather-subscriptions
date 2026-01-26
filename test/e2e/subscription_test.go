package e2e

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gorm.io/gorm"

	"weather-subscriptions/internal/db/models"
)

//nolint:gocognit,gocycle,gocyclo,cyclop
func TestSubscriptions(t *testing.T) {
	webApp, state := newTestApp()
	app := adaptor.FiberApp(webApp)

	suffix := time.Now().UnixNano()
	cityName := fmt.Sprintf("warsaw-%d", suffix)

	city := &models.City{
		ID:            fmt.Sprintf("city-%d", suffix),
		Name:          cityName,
		Longitude:     21.0122,
		Latitude:      52.2297,
		GooglePlaceID: fmt.Sprintf("place-%d", suffix),
	}
	if err := state.SaveCity(city); err != nil {
		t.Fatalf("expected city to be saved: %v", err)
	}

	firstEmail := fmt.Sprintf("user-%d@example.com", suffix)
	secondEmail := fmt.Sprintf("user-%d@example.com", suffix+1)

	var user *models.User
	var confirmToken *models.Token
	var unsubToken *models.Token
	var subscription *models.Subscription
	var secondUser *models.User
	var secondConfirmToken *models.Token
	var secondUnsubToken *models.Token

	t.Run("subscribe creates user and tokens", func(t *testing.T) {
		request := map[string]string{
			"email":     firstEmail,
			"city":      cityName,
			"frequency": models.DAILY,
		}
		response := sendJSONRequest(t, app, http.MethodPost, "/subscribe", request)
		if response.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.StatusCode)
		}

		var err error
		user, err = state.GetUserByEmail(firstEmail)
		if err != nil {
			t.Fatalf("expected user to be stored: %v", err)
		}
		if user == nil {
			t.Fatalf("expected user to be created")
		}

		confirmToken, err = state.GetTokenByUserIDAndType(user.ID, models.Sub)
		if err != nil {
			t.Fatalf("expected confirmation token: %v", err)
		}
		if confirmToken == nil {
			t.Fatalf("expected confirmation token to exist")
		}

		unsubToken, err = state.GetTokenByUserIDAndType(user.ID, models.Unsub)
		if err != nil {
			t.Fatalf("expected unsubscribe token: %v", err)
		}
		if unsubToken == nil {
			t.Fatalf("expected unsubscribe token to exist")
		}
	})

	t.Run("confirm subscription with valid token", func(t *testing.T) {
		confirmResponse := sendRequest(app, http.MethodGet, "/confirm/"+confirmToken.Token, nil)
		if confirmResponse.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", confirmResponse.StatusCode)
		}

		var err error
		subscription, err = state.GetSubscription(user.ID)
		if err != nil {
			t.Fatalf("expected subscription: %v", err)
		}
		if subscription == nil {
			t.Fatalf("expected subscription to be created")
		}
		if subscription.Frequency != models.DAILY {
			t.Fatalf("expected frequency %q, got %q", models.DAILY, subscription.Frequency)
		}

		_, err = state.GetToken(confirmToken.Token)
		if err == nil {
			t.Fatalf("expected confirmation token to be removed")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("unexpected error looking up token: %v", err)
		}
	})

	t.Run("confirm subscription with invalid token", func(t *testing.T) {
		invalidConfirmResponse := sendRequest(app, http.MethodGet, "/confirm/invalid", nil)
		if invalidConfirmResponse.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", invalidConfirmResponse.StatusCode)
		}
	})

	t.Run("confirm subscription with expired token", func(t *testing.T) {
		secondRequest := map[string]string{
			"email":     secondEmail,
			"city":      cityName,
			"frequency": models.HOURLY,
		}
		secondResponse := sendJSONRequest(t, app, http.MethodPost, "/subscribe", secondRequest)
		if secondResponse.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", secondResponse.StatusCode)
		}

		var err error
		secondUser, err = state.GetUserByEmail(secondEmail)
		if err != nil {
			t.Fatalf("expected second user to be stored: %v", err)
		}

		secondConfirmToken, err = state.GetTokenByUserIDAndType(secondUser.ID, models.Sub)
		if err != nil {
			t.Fatalf("expected second confirmation token: %v", err)
		}
		secondConfirmToken.ExpiryAt = time.Now().Add(-time.Hour)
		secondConfirmToken.CreatedAt = time.Now().Add(-time.Hour)
		secondConfirmToken.UpdatedAt = time.Now().Add(-time.Hour)
		if err := state.SaveToken(secondConfirmToken); err != nil {
			t.Fatalf("expected expired token to be saved: %v", err)
		}

		expiredConfirmResponse := sendRequest(app, http.MethodGet, "/confirm/"+secondConfirmToken.Token, nil)
		if expiredConfirmResponse.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 400, got %d", expiredConfirmResponse.StatusCode)
		}
	})

	t.Run("unsubscribe with valid token", func(t *testing.T) {
		unsubscribeResponse := sendRequest(app, http.MethodGet, "/unsubscribe/"+unsubToken.Token, nil)
		if unsubscribeResponse.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", unsubscribeResponse.StatusCode)
		}

		removedUser, err := state.GetUser(user.ID)
		if err == nil || removedUser != nil {
			t.Fatalf("expected user to be removed")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("unexpected error looking up user: %v", err)
		}

		_, err = state.GetSubscription(user.ID)
		if err == nil {
			t.Fatalf("expected subscription to be removed")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("unexpected error looking up subscription: %v", err)
		}

		_, err = state.GetToken(unsubToken.Token)
		if err == nil {
			t.Fatalf("expected unsubscribe token to be removed")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("unexpected error looking up unsubscribe token: %v", err)
		}
	})

	t.Run("unsubscribe with invalid token", func(t *testing.T) {
		invalidUnsubResponse := sendRequest(app, http.MethodGet, "/unsubscribe/invalid", nil)
		if invalidUnsubResponse.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", invalidUnsubResponse.StatusCode)
		}
	})

	t.Run("unsubscribe with expired token", func(t *testing.T) {
		var err error
		secondUnsubToken, err = state.GetTokenByUserIDAndType(secondUser.ID, models.Unsub)
		if err != nil {
			t.Fatalf("expected second unsubscribe token: %v", err)
		}
		secondUnsubToken.ExpiryAt = time.Now().Add(-time.Hour)
		secondUnsubToken.CreatedAt = time.Now().Add(-time.Hour)
		secondUnsubToken.UpdatedAt = time.Now().Add(-time.Hour)
		if err := state.SaveToken(secondUnsubToken); err != nil {
			t.Fatalf("expected expired unsub token to be saved: %v", err)
		}

		expiredUnsubResponse := sendRequest(
			app,
			http.MethodGet,
			fmt.Sprintf("/unsubscribe/%s", secondUnsubToken.Token),
			nil,
		)
		if expiredUnsubResponse.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", expiredUnsubResponse.StatusCode)
		}
	})
}

func sendJSONRequest(
	t *testing.T,
	app http.Handler,
	method string,
	path string,
	payload any,
) *http.Response {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	return sendRequest(app, method, path, bytes.NewReader(body))
}

func sendRequest(
	app http.Handler,
	method string,
	path string,
	body io.Reader,
) *http.Response {
	request := httptest.NewRequest(method, path, body)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response.Result()
}
