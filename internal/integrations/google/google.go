package google

import (
	"context"
	"github.com/google/uuid"
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
		return nil, err
	}

	return weather, nil
}

func (g *Google) GetCity(ctx context.Context, cityName string) (*models.City, error) {
	mapsClient, err := maps.NewClient(maps.WithAPIKey(g.cfg.GoogleMapsApiKey))
	if err != nil {
		return nil, err
	}
	return getCity(ctx, mapsClient, cityName)
}

//func (g *Google) GetWeatherByCityID(ctx context.Context, state state.Stateful, cityID string) (*models.Weather, error) {
//	city, err := state.GetCityByID(cityID)
//	if err != nil {
//		return nil, err
//	}
//
//	weather, err := g.fetchWeatherForCity(ctx, city)
//	if err != nil {
//		return nil, err
//	}
//
//	err = state.SaveWeather(weather)
//	if err != nil {
//		return nil, err
//	}
//
//	return weather, nil
//}

func getCity(ctx context.Context, client *maps.Client, cityName string) (*models.City, error) {
	cityInfo, err := fetchCityInfo(ctx, client, cityName)
	if err != nil {
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
