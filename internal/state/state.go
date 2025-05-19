package state

import (
	"gorm.io/gorm"
	"strings"
	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/state/resolvers"
)

type Stateful interface {
	GetUser(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetCity(name string) (*models.City, error)
	GetCityByID(id string) (*models.City, error)
	GetWeather(cityID string) (*models.Weather, error)
	GetToken(tokens string) (*models.Token, error)
	GetUnsubToken(userID string) (*models.Token, error)
	GetSubToken(userID string) (*models.Token, error)
	GetSubscription(userID string) (*models.Subscription, error)
	GetSubscriptions(subscriptionType models.SubscriptionType) ([]*models.Subscription, error)
	SaveWeather(weather *models.Weather) error
	SaveCity(city *models.City) error
	SaveUser(user *models.User) error
	SaveToken(token *models.Token) error
	SaveSubscription(subscription *models.Subscription) error
	RemoveSubscription(subscription *models.Subscription) error
	RemoveToken(token *models.Token) error
	RemoveUser(user *models.User) error
}

type State struct {
	resolver      resolvers.Resolver
	user          map[string]*models.User
	cities        map[string]*models.City
	cityIDMap     map[string]*models.City
	weather       map[string]*models.Weather
	tokens        map[string]*models.Token
	subscriptions map[string]*models.Subscription
}

func (s *State) GetUser(id string) (*models.User, error) {
	user, ok := s.user[id]
	if !ok {
		foundUser, err := s.resolver.UserByID(id)
		if err != nil {
			return nil, err
		}
		user = foundUser
	}
	s.user[id] = user

	return user, nil
}

func (s *State) GetUserByEmail(email string) (*models.User, error) {
	user, ok := s.user[email]
	if !ok {
		foundUser, err := s.resolver.UserByEmail(email)
		if err != nil {
			return nil, err
		}
		user = foundUser
	}

	s.user[email] = user
	return user, nil
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

func (s *State) GetToken(token string) (*models.Token, error) {
	userToken, ok := s.tokens[token]
	if !ok {
		foundToken, err := s.resolver.Token(token)
		if err != nil {
			return nil, err
		}
		userToken = foundToken
	}
	s.tokens[token] = userToken

	return userToken, nil
}

func (s *State) GetSubToken(userID string) (*models.Token, error) {
	token, err := s.resolver.SubToken(userID)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *State) GetUnsubToken(userID string) (*models.Token, error) {
	token, err := s.resolver.UnsubToken(userID)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *State) GetSubscription(userID string) (*models.Subscription, error) {
	subscription, ok := s.subscriptions[userID]
	if !ok {
		foundSubscription, err := s.resolver.Subscription(userID)
		if err != nil {
			return nil, err
		}
		subscription = foundSubscription
	}
	s.subscriptions[userID] = subscription

	return subscription, nil
}

func (s *State) GetSubscriptions(subscriptionType models.SubscriptionType) (subs []*models.Subscription, err error) {
	subscriptions, err := s.resolver.Subscriptions(subscriptionType)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
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

func (s *State) SaveWeather(weather *models.Weather) error {
	err := s.resolver.Save(weather)
	if err != nil {
		return err
	}
	s.weather[weather.CityID] = weather

	return nil
}

func (s *State) SaveCity(city *models.City) error {
	err := s.resolver.Save(city)
	if err != nil {
		return err
	}
	s.cities[strings.ToLower(city.Name)] = city

	return nil
}

func (s *State) SaveUser(user *models.User) error {
	err := s.resolver.Save(user)
	if err != nil {
		return err
	}
	s.user[user.Email] = user
	s.user[user.ID] = user

	return nil
}

func (s *State) SaveToken(token *models.Token) error {
	err := s.resolver.Save(token)
	if err != nil {
		return err
	}
	s.tokens[token.Token] = token

	return nil
}

func (s *State) SaveSubscription(subscription *models.Subscription) error {
	err := s.resolver.Save(subscription)
	if err != nil {
		return err
	}
	s.subscriptions[subscription.UserID] = subscription

	return nil
}

func (s *State) RemoveSubscription(subscription *models.Subscription) error {
	err := s.resolver.Remove(subscription)
	if err != nil {
		return err
	}

	delete(s.subscriptions, subscription.UserID)
	delete(s.user, subscription.UserID)

	return nil
}

func (s *State) RemoveUser(user *models.User) error {
	err := s.resolver.Remove(user)
	if err != nil {
		return err
	}

	delete(s.user, user.ID)
	delete(s.subscriptions, user.ID)

	return nil
}

func (s *State) RemoveToken(token *models.Token) error {
	err := s.resolver.Remove(token)
	if err != nil {
		return err
	}

	delete(s.tokens, token.Token)

	return nil
}

func NewState(db *gorm.DB) Stateful {
	resolver := resolvers.New(db)
	return &State{
		resolver:      resolver,
		user:          make(map[string]*models.User),
		cities:        make(map[string]*models.City),
		cityIDMap:     make(map[string]*models.City),
		weather:       make(map[string]*models.Weather),
		tokens:        make(map[string]*models.Token),
		subscriptions: make(map[string]*models.Subscription),
	}
}
