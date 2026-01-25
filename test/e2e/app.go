package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	handlers "weather-subscriptions/api/handlers/subscription"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/mail"
	"weather-subscriptions/internal/state"
)

// testMailStub is a no-op mailer for e2e tests.
type testMailStub struct{}

func (m *testMailStub) Send(_ mail.Message) error {
	return nil
}

// testMapsStub returns a deterministic city for e2e tests.
type testMapsStub struct{}

func (m *testMapsStub) GetWeather(_ context.Context, _ *models.City) (*models.Weather, error) {
	return nil, fmt.Errorf("not implemented in tests")
}

func (m *testMapsStub) GetCity(_ context.Context, cityName string) (*models.City, error) {
	return &models.City{
		ID:            uuid.Must(uuid.NewV7()).String(),
		Name:          cityName,
		Latitude:      52.2297,
		Longitude:     21.0122,
		GooglePlaceID: "test-place-id",
	}, nil
}

type testApp struct {
	App    *fiber.App
	DB     *gorm.DB
	State  *state.State
	Config *config.Config
}

// newTestApp sets up Fiber app, DB and state with stubbed integrations and mailer.
func newTestApp(t *testing.T) *testApp {
	t.Helper()

	dns := os.Getenv("DNS")
	if dns == "" {
		t.Fatal("DNS environment variable must be set for e2e tests")
	}

	cfg := &config.Config{
		DNS:         dns,
		Port:        "3000",
		FrontendURL: "http://example.com",
	}

	database, err := db.Connect(cfg)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	st := state.NewState(database)
	mailerStub := &testMailStub{}
	mapsStub := &testMapsStub{}

	subHandler := handlers.NewSubscriptionHandler(cfg, st, mailerStub, mapsStub)

	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Post("/subscribe", subHandler.HandleSubscribe)
	app.Get("/confirm/:token", subHandler.HandleConfirmSubscription)
	app.Get("/unsubscribe/:token", subHandler.HandleUnsubscribe)

	return &testApp{
		App:    app,
		DB:     database,
		State:  st,
		Config: cfg,
	}
}

// resetTables truncates core tables to keep tests isolated.
func resetTables(t *testing.T, dbConn *gorm.DB) {
	t.Helper()
	tables := []string{"advises", "subscriptions", "tokens", "users", "cities", "weathers"}
	for _, tbl := range tables {
		if err := dbConn.Exec(fmt.Sprintf("DELETE FROM %s", tbl)).Error; err != nil {
			t.Fatalf("failed to clean table %s: %v", tbl, err)
		}
	}
}

// seedExpiredToken inserts an expired token for the given user.
func seedExpiredToken(t *testing.T, dbConn *gorm.DB, userID string, tokenType models.TokenType) string {
	t.Helper()
	token := &models.Token{
		Token:            uuid.Must(uuid.NewV7()).String(),
		Type:             tokenType,
		SubscriptionType: nil,
		UserID:           userID,
		CreatedAt:        time.Now().Add(-48 * time.Hour),
		UpdatedAt:        time.Now().Add(-48 * time.Hour),
		ExpiryAt:         time.Now().Add(-24 * time.Hour),
	}
	if err := dbConn.Create(token).Error; err != nil {
		t.Fatalf("failed to seed expired token: %v", err)
	}
	return token.Token
}
