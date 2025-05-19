package google

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/integrations"
)

type Google struct {
	cfg *config.Config
}

func New(cfg *config.Config) integrations.MapsIntegration {
	return &Google{
		cfg: cfg,
	}
}

func (g *Google) GetWeather(ctx context.Context, city *models.City) (*models.Weather, error) {
	weather, err := g.fetchWeatherForCity(ctx, city)
	if err != nil {
		zap.L().Error("failed to fetch weather", zap.Error(err))
		return nil, err
	}

	return weather, nil
}

func (g *Google) GetCity(ctx context.Context, cityName string) (*models.City, error) {
	mapsClient, err := maps.NewClient(maps.WithAPIKey(g.cfg.GoogleMapsApiKey))
	if err != nil {
		zap.L().Error("failed to create maps client", zap.Error(err))
		return nil, err
	}
	return getCity(ctx, mapsClient, cityName)
}

func getCity(ctx context.Context, client *maps.Client, cityName string) (*models.City, error) {
	cityInfo, err := fetchCityInfo(ctx, client, cityName)
	if err != nil {
		zap.L().Error("failed to fetch city info", zap.Error(err))
		return nil, err
	}

	city := &models.City{
		ID:            uuid.Must(uuid.NewV7()).String(),
		Name:          cityName,
		Latitude:      cityInfo.Latitude,
		Longitude:     cityInfo.Longitude,
		GooglePlaceID: cityInfo.GooglePlaceID,
	}

	return city, nil
}
