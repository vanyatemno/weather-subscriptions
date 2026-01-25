package state

import (
	"strings"
	"weather-subscriptions/internal/db/models"
)

type CitiesState interface {
	GetCity(name string) (*models.City, error)
	GetCityByID(id string) (*models.City, error)
	SaveCity(city *models.City) error
}

func (s *State) GetCity(name string) (*models.City, error) {
	city, ok := s.cities[strings.ToLower(name)]

	if !ok {
		foundCity, err := s.resolver.City(name)
		if err != nil {
			return nil, err
		}

		city = foundCity
	}
	s.cities[name] = city

	return city, nil
}

// todo: remove
func (s *State) GetCityByID(id string) (*models.City, error) {
	city, ok := s.cityIDMap[id]
	if !ok {
		foundCity, err := s.resolver.CityByID(id)
		if err != nil {
			return nil, err
		}
		city = foundCity
	}

	return city, nil
}

func (s *State) SaveCity(city *models.City) error {
	err := s.resolver.Save(city)
	if err != nil {
		return err
	}
	s.cities[strings.ToLower(city.Name)] = city

	return nil
}
