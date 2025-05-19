package models

import "time"

type Weather struct {
	ID          string    `gorm:"primaryKey;default:uuid_generate_v4()"`
	Time        time.Time `gorm:"not null"`
	Temperature float64   `gorm:"not null"`
	Humidity    int       `gorm:"not null"`
	Description string    `gorm:"not null"`
	CityID      string    `gorm:"not null"`
	City        City      `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
