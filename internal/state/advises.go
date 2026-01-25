package state

import "weather-subscriptions/internal/db/models"

type AdvisesState interface {
	GetAdvises(weatherID string) ([]*models.Advise, error)
	SaveAdvises(advises []*models.Advise) error
	RemoveAdvise(advise *models.Advise) error
}

func (s *State) GetAdvises(weatherID string) ([]*models.Advise, error) {
	advises, ok := s.advises[weatherID]
	if !ok {
		advises, err := s.resolver.AdviseByWeatherID(weatherID)
		if err != nil {
			return nil, err
		}
		return advises, nil
	}

	return advises, nil
}

func (s *State) SaveAdvises(advises []*models.Advise) error {
	for _, advise := range advises {
		err := s.resolver.Save(advise)
		if err != nil {
			return err
		}
	}

	s.advises[advises[0].WeatherID] = advises

	return nil
}

func (s *State) RemoveAdvise(advise *models.Advise) error {
	err := s.resolver.Remove(advise)
	if err != nil {
		return err
	}

	advises, ok := s.advises[advise.WeatherID]
	if ok {
		var filteredAdvises []*models.Advise
		for i := range advises {
			if advises[i].ID != advise.ID {
				filteredAdvises = append(filteredAdvises, advises[i])
			}
		}
		s.advises[advise.WeatherID] = filteredAdvises
	}

	return nil
}
