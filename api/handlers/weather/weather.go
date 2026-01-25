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

type WeatherResponse struct {
	Temperature float64 `json:"temperature" example:"21.5"`
	Humidity    float64 `json:"humidity" example:"64"`
	Description string  `json:"description" example:"partly cloudy"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"city name is required"`
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

// GetWeather godoc
// @Summary Get current weather for a city
// @Description Returns the current weather forecast for the specified city.
// @Tags weather
// @Param city query string true "City name"
// @Success 200 {object} WeatherResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /weather [get]
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
