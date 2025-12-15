package templates

import (
	"fmt"
	"strconv"
	"weather-subscriptions/internal/db/models"
)

const (
	unsubscribeLinkTemplate = "%s/unsubscribe/%s"
	subscribeLinkTemplate   = "%s/confirm/%s"
)

func GetWeatherEmailBody(
	weather *models.Weather,
	frontendURL, code string,
) string {
	return fmt.Sprintf(
		weatherEmailTemplate,
		strconv.FormatFloat(weather.Temperature, 'f', -1, 64),
		strconv.Itoa(weather.Humidity),
		weather.Description,
		fmt.Sprintf(unsubscribeLinkTemplate, frontendURL, code),
	)
}

func GetVerificationEmailTemplate(frontendURL, code string) string {
	subscribeLink := fmt.Sprintf(subscribeLinkTemplate, frontendURL, code)
	return fmt.Sprintf(
		verificationEmailTemplate,
		subscribeLink,
	)
}
