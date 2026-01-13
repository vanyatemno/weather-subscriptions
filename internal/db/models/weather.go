package models

import "time"

type Weather struct {
	ID     string `gorm:"primaryKey;default:uuid_generate_v4()"`
	CityID string

	Time        time.Time
	Temperature float64
	Humidity    int
	Description string

	City City `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
