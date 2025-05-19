package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/state"
)

type WeatherHandler struct {
	googleInt integrations.MapsIntegration
	state     state.Stateful
}

func NewWeatherHandler(googleInt integrations.MapsIntegration, state state.Stateful) *WeatherHandler {
	return &WeatherHandler{googleInt: googleInt, state: state}
}

func (wh *WeatherHandler) GetWeather(c *fiber.Ctx) error {
	cityName := c.Query("city")
	if cityName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "city name is required"})
	}
	cityName = slug.Make(cityName)

	city, err := wh.state.GetCity(cityName)
	if err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		city, err = wh.googleInt.GetCity(c.Context(), cityName)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		err = wh.state.SaveCity(city)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	} else if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	weather, err := wh.googleInt.GetWeather(c.Context(), city)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err = wh.state.SaveWeather(weather)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"temperature": weather.Temperature,
		"humidity":    weather.Humidity,
		"description": weather.Description,
	})
}
