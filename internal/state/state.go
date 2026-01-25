package state

import (
	"gorm.io/gorm"

	"weather-subscriptions/internal/db/models"
	"weather-subscriptions/internal/state/resolvers"
)

type State struct {
	resolver      resolvers.Resolver
	user          map[string]*models.User
	cities        map[string]*models.City
	cityIDMap     map[string]*models.City
	weather       map[string]*models.Weather
	tokens        map[string]*models.Token
	subscriptions map[string]*models.Subscription
	advises       map[string][]*models.Advise
}

func NewState(db *gorm.DB) *State {
	resolver := resolvers.New(db)
	return &State{
		resolver:      resolver,
		user:          make(map[string]*models.User),
		cities:        make(map[string]*models.City),
		cityIDMap:     make(map[string]*models.City),
		weather:       make(map[string]*models.Weather),
		tokens:        make(map[string]*models.Token),
		subscriptions: make(map[string]*models.Subscription),
		advises:       make(map[string][]*models.Advise),
	}
}
