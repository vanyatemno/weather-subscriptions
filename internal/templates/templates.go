package templates

import (
	"fmt"
	"strconv"
	"weather-subscriptions/internal/db/models"
)

func GetWeatherEmailBody(weather *models.Weather, code string) string {
	return fmt.Sprintf(
		weatherEmailTemplate,
		strconv.FormatFloat(weather.Temperature, 'f', -1, 64),
		strconv.Itoa(weather.Humidity),
		weather.Description,
		code,
	)
}

func GetVerificationEmailTemplate(code string) string {
	return fmt.Sprintf(verificationEmailTemplate, code)
}
