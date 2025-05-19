package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"net/url"
	"time"
	"weather-subscriptions/internal/db/models"
)

const weatherURL = "https://weather.googleapis.com/v1/currentConditions:lookup"

func (g *Google) fetchWeatherForCity(ctx context.Context, city *models.City) (*models.Weather, error) {
	client := resty.New()
	query := g.getQuery(city)

	var result WeatherResponse
	req, err := client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(weatherURL + "?" + query)
	if err != nil {
		return nil, err
	}
	if !req.IsSuccess() {
		fmt.Println(string(req.Body()))
		return nil, errors.New("could not fetch weather for city: " + city.Name)
	}

	return &models.Weather{
		ID:          uuid.Must(uuid.NewV7()).String(),
		Time:        time.Now(),
		Temperature: result.Temperature.Degrees,
		Humidity:    result.RelativeHumidity,
		Description: result.WeatherCondition.Description.Text,
		CityID:      city.ID,
		City:        *city,
	}, nil
}

func (g *Google) getQuery(city *models.City) string {
	query := url.Values{}
	cityCoordinates := city.GetStringCoordinates()
	query.Set("key", g.cfg.GoogleMapsApiKey)
	query.Set("location.latitude", cityCoordinates.Lat)
	query.Set("location.longitude", cityCoordinates.Long)

	return query.Encode()
}
