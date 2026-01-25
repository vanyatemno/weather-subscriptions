package advises

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/integrations"
	"weather-subscriptions/internal/state"
	"weather-subscriptions/internal/util"
)

const prompt = `Role & Objective: You are an expert Local Guide and Travel Concierge.
Your task is to analyze TripAdvisor listings for a specific location and provide tailored recommendations based strictly on the current weather conditions.
You must adapt your search strategy to ensure user comfort and enjoyment.

Input Variables:

Location: %s

Current Weather: %s

Instructions:

Step 1: Analyze the Weather Context Determine if the weather is "Good" or "Bad" for standard tourism based on the input:

"Bad" Weather (Rain, Snow, Extreme Cold/Heat, High Wind): Focus on Shelter & Comfort. The goal is "Cozy/Indoor."

"Good" Weather (Sunny, Mild, Clear Skies): Focus on Exploration & Scenery. The goal is "Active/Outdoor."

Step 2: Apply Search & Filter Strategy (TripAdvisor Simulation) Based on the weather analysis, simulate a search for the following criteria:

IF WEATHER IS BAD (Seek "The Great Indoors"):

Keywords: Search for terms like "cozy," "fireplace," "hidden gem," "atmosphere," "books," "cat café," "underground," "museum," "spa."

Categories:

Food & Drink: Coffee shops, tearooms, historic pubs/bars, restaurants with long tasting menus.

Attractions: Art galleries, interactive museums, escape rooms, bowling, covered markets.

Vibe Check: Look for reviews mentioning "great place to escape the rain" or "warm atmosphere."

IF WEATHER IS GOOD (Seek "The Great Outdoors"):

Keywords: Search for terms like "terrace," "rooftop," "garden," "view," "panoramic," "outdoor seating," "walking tour," "park."

Categories:

Food & Drink: Rooftop bars, beer gardens, riverside dining, street food markets.

Attractions: Botanical gardens, open-air museums, historical walking areas, observation decks, hiking trails, boat rentals.

Vibe Check: Look for reviews mentioning "beautiful sunset," "best view in the city," or "lovely walk."

Step 3: Generate Recommendations Provide a curated list of 3-5 distinct options. Do not include any links in place description text.`

type AdvisesService struct {
	advisesState state.AdvisesState
	provider     integrations.Generator
}

func NewAdvisesService(advisesState state.AdvisesState, provider integrations.Generator) *AdvisesService {
	return &AdvisesService{
		advisesState: advisesState,
		provider:     provider,
	}
}

func (a *AdvisesService) GetAdvise(
	ctx context.Context,
	weather *models.Weather,
) (*Advise, error) {
	var res Advise

	advises, err := a.advisesState.GetAdvises(weather.ID)
	if err != nil {
		return nil, err
	}

	if len(advises) != 0 {
		for i := range advises {
			res.Places = append(res.Places, Place{
				Name:        advises[i].Name,
				Link:        advises[i].Link,
				Description: advises[i].Description,
			})
		}
	}

	err = a.provider.GenerateStructuredResponse(ctx, buildPrompt(weather), &res)
	if err != nil {
		zap.L().Error("Failed to generate structured response", zap.Error(err))
		return nil, err
	}
	// remove links provided by AI as a default
	cleanupResponse(&res)

	err = a.saveAdvise(&res, weather)
	if err != nil {
		zap.L().Error("Failed to save structured response", zap.Error(err))
		return nil, err
	}

	return &res, nil
}

func buildPrompt(weather *models.Weather) string {
	return fmt.Sprintf(
		prompt,
		weather.City.Name,
		fmt.Sprintf("%f°C %s", weather.Temperature, weather.Description),
	)
}

func (a *AdvisesService) saveAdvise(advise *Advise, weather *models.Weather) error {
	var advisesModels []*models.Advise
	for _, place := range advise.Places {
		advisesModels = append(advisesModels, &models.Advise{
			WeatherID:   weather.ID,
			Name:        place.Name,
			Description: place.Description,
			Link:        place.Link,
		})
	}

	return a.advisesState.SaveAdvises(advisesModels)
}

// cleanupResponse - removes all source links provided in AI response
func cleanupResponse(res *Advise) {
	for i := range res.Places {
		res.Places[i].Description = util.RemoveAiLinks(res.Places[i].Description)
	}
}
