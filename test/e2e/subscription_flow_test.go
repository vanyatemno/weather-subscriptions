package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"weather-subscriptions/internal/db/models"
)

type subscribeRequest struct {
	Email     string `json:"email"`
	City      string `json:"city"`
	Frequency string `json:"frequency"`
}

func doRequest(t *testing.T, server *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req, err := http.NewRequest(method, server.URL+path, &buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func findTokenByType(t *testing.T, dbConn *gorm.DB, tokenType models.TokenType) *models.Token {
	t.Helper()
	var token models.Token
	if err := dbConn.Where("type = ?", tokenType).First(&token).Error; err != nil {
		t.Fatalf("failed to find token type %s: %v", tokenType, err)
	}
	return &token
}

func TestSubscriptionFlow_HappyPath(t *testing.T) {
	app := newTestApp(t)
	resetTables(t, app.DB)

	srv := httptest.NewServer(adaptor.FiberApp(app.App))
	defer srv.Close()

	// Subscribe
	resp := doRequest(t, srv, http.MethodPost, "/subscribe", subscribeRequest{
		Email:     "user@example.com",
		City:      "warsaw",
		Frequency: string(models.DAILY),
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("subscribe status = %d", resp.StatusCode)
	}

	// confirm using generated sub token
	subToken := findTokenByType(t, app.DB, models.Sub)
	resp = doRequest(t, srv, http.MethodGet, "/confirm/"+subToken.Token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("confirm status = %d", resp.StatusCode)
	}

	// unsubscribe using generated unsub token
	unsubToken := findTokenByType(t, app.DB, models.Unsub)
	resp = doRequest(t, srv, http.MethodGet, "/unsubscribe/"+unsubToken.Token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unsubscribe status = %d", resp.StatusCode)
	}
}

func TestSubscriptionFlow_InvalidToken(t *testing.T) {
	app := newTestApp(t)
	resetTables(t, app.DB)

	srv := httptest.NewServer(adaptor.FiberApp(app.App))
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/confirm/"+uuid.NewString(), nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for invalid token, got %d", resp.StatusCode)
	}
}

func TestSubscriptionFlow_ExpiredToken(t *testing.T) {
	app := newTestApp(t)
	resetTables(t, app.DB)

	// seed user and expired sub token
	user := &models.User{
		ID:     uuid.Must(uuid.NewV7()).String(),
		Email:  "exp@example.com",
		CityID: uuid.Must(uuid.NewV7()).String(),
	}
	city := &models.City{
		ID:            user.CityID,
		Name:          "warsaw",
		Latitude:      52.2297,
		Longitude:     21.0122,
		GooglePlaceID: "test-place",
	}
	if err := app.DB.Create(city).Error; err != nil {
		t.Fatalf("seed city: %v", err)
	}
	if err := app.DB.Create(user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	expired := seedExpiredToken(t, app.DB, user.ID, models.Sub)

	srv := httptest.NewServer(adaptor.FiberApp(app.App))
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/confirm/"+expired, nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for expired token, got %d", resp.StatusCode)
	}
}

func TestSubscriptionFlow_UnsubscribeInvalid(t *testing.T) {
	app := newTestApp(t)
	resetTables(t, app.DB)

	srv := httptest.NewServer(adaptor.FiberApp(app.App))
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/unsubscribe/"+uuid.NewString(), nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for invalid unsubscribe token, got %d", resp.StatusCode)
	}
}
