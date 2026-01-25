package templates

import (
	"fmt"
	"strconv"
	"strings"
	"weather-subscriptions/internal/advises"
	"weather-subscriptions/internal/db/models"
)

const (
	unsubscribeLinkTemplate = "%s/unsubscribe/%s"
	subscribeLinkTemplate   = "%s/confirm/%s"
	recommendationTemplate  = `<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f8f9fa; border-radius: 8px; margin: 10px 0; padding: 15px; border-left: 4px solid #3498db;">
    <tr>
        <td>
            <h3 style="margin: 0 0 8px 0; color: #2c3e50; font-size: 18px;">%s</h3>
            <p style="margin: 0 0 10px 0; color: #34495e; font-size: 14px; line-height: 1.5;">%s</p>
            <a href="%s" style="display: inline-block; background-color: #3498db; color: white; text-decoration: none; padding: 8px 16px; border-radius: 5px; font-size: 14px;">View on TripAdvisor</a>
        </td>
    </tr>
</table>`
)

func GetWeatherEmailBody(
	weather *models.Weather,
	frontendURL, code string,
	advises *advises.Advise,
) string {
	return fmt.Sprintf(
		weatherEmailTemplate,
		strconv.FormatFloat(weather.Temperature, 'f', -1, 64),
		strconv.Itoa(weather.Humidity),
		weather.Description,
		fmt.Sprintf(unsubscribeLinkTemplate, frontendURL, code),
		buildRecommendations(advises),
	)
}

func buildRecommendations(advise *advises.Advise) string {
	var recommendations []string
	for i := range advise.Places {
		recommendation := fmt.Sprintf(
			recommendationTemplate,
			advise.Places[i].Name,
			advise.Places[i].Description,
			advise.Places[i].Link,
		)
		recommendations = append(recommendations, recommendation)
	}

	return strings.Join(recommendations, "\n")
}

func GetVerificationEmailTemplate(frontendURL, code string) string {
	subscribeLink := fmt.Sprintf(subscribeLinkTemplate, frontendURL, code)
	return fmt.Sprintf(
		verificationEmailTemplate,
		subscribeLink,
	)
}
