package integrations

import (
	"context"
	"weather-subscriptions/internal/db/models"
)

// MapsIntegration interface to all integrations which fetch data about city coordinates or current weather
type MapsIntegration interface {
	GetWeather(ctx context.Context, city *models.City) (*models.Weather, error)
	GetCity(ctx context.Context, cityName string) (*models.City, error)
}

type Generator interface {
	GeneratePlainResponse(ctx context.Context, prompt string) (string, error)
	GenerateStructuredResponse(ctx context.Context, prompt string, dest any) error
}
