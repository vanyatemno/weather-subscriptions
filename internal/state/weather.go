package state

import "weather-subscriptions/internal/db/models"

type WeatherState interface {
	GetWeather(cityID string) (*models.Weather, error)
	SaveWeather(weather *models.Weather) error
}

func (s *State) GetWeather(cityID string) (*models.Weather, error) {
	weather, ok := s.weather[cityID]
	if !ok {
		foundWeather, err := s.resolver.WeatherByCityID(cityID)
		if err != nil {
			return nil, err
		}
		weather = foundWeather
	}
	s.weather[cityID] = weather

	return weather, nil
}

func (s *State) SaveWeather(weather *models.Weather) error {
	err := s.resolver.Save(weather)
	if err != nil {
		return err
	}
	s.weather[weather.CityID] = weather

	return nil
}
