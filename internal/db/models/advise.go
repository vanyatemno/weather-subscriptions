package models

import "time"

type Advise struct {
	ID        string `gorm:"primaryKey;default:uuid_generate_v4()"`
	WeatherID string

	Name        string
	Description string
	Link        string

	Weather Weather

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
