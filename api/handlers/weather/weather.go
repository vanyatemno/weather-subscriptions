package handlers

import (
	"errors"

	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/state"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type WeatherHandler struct {
	googleInt    integrations.MapsIntegration
	weatherState state.WeatherState
	citiesState  state.CitiesState
}

func NewWeatherHandler(
	googleInt integrations.MapsIntegration,
	weatherState state.WeatherState,
	citiesState state.CitiesState,
) *WeatherHandler {
	return &WeatherHandler{
		googleInt:    googleInt,
		weatherState: weatherState,
		citiesState:  citiesState,
	}
}

func (wh *WeatherHandler) GetWeather(c *fiber.Ctx) error {
	cityName := c.Query("city")
	if cityName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "city name is required"})
	}
	cityName = slug.Make(cityName)

	city, err := wh.citiesState.GetCity(cityName)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		city, err = wh.googleInt.GetCity(c.Context(), cityName)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		err = wh.citiesState.SaveCity(city)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	} else if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	weather, err := wh.googleInt.GetWeather(c.Context(), city)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = wh.weatherState.SaveWeather(weather)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"temperature": weather.Temperature,
		"humidity":    weather.Humidity,
		"description": weather.Description,
	})
}
