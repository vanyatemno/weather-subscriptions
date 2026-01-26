package e2e

import (
	"os"
	"strconv"
	"weather-subscriptions/api/routes"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/test/mock"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func newTestApp() (*fiber.App, *state.State) {
	cfg := newTestConfig()
	database, err := db.Connect(cfg)
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	mailerService := mock.NewMailerServiceMock()
	set := state.NewState(database)

	webApp := fiber.New()
	webApp.Use(
		cors.New(
			cors.Config{
				AllowOrigins: "*",
			},
		),
	)
	routes.New(cfg, set, mailerService).Setup(webApp)

	return webApp, set
}

func newTestConfig() *config.Config {
	mailerPort, err := strconv.Atoi(os.Getenv("MAILER_PORT"))
	if err != nil {
		mailerPort = 465
	}

	return &config.Config{
		DNS: "postgresql://postgres:Password1@localhost:5432/weather-subscriptions",
		//DNS:              os.Getenv("DNS"),
		Port:             "3000",
		FrontendURL:      "http://example.com",
		GoogleMapsAPIKey: "test-key",
		Mailer: config.Mailer{
			Host:     os.Getenv("MAILER_HOST"),
			Port:     mailerPort,
			Username: os.Getenv("MAILER_USERNAME"),
			From:     os.Getenv("MAILER_FROM"),
			SMTP:     os.Getenv("MAILER_SMTP"),
			Password: os.Getenv("MAILER_PASSWORD"),
		},
		OpenAI: config.OpenAI{
			OpenrouterAPIKey: "test-key",
		},
	}
}
