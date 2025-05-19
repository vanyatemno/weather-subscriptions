package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db/models"
)

func Connect(config *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(config.DNS), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = database.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return nil, err
	}
	err = database.AutoMigrate(
		&models.City{},
		&models.User{},
		&models.Token{},
		&models.Weather{},
		&models.Subscription{},
	)
	if err != nil {
		return nil, err
	}

	return database, nil
}
