package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;default:uuid_generate_v4()"`
	Email     string    `gorm:"not null;unique"`
	CityID    string    `gorm:"not null;foreignKey:CityID"`
	City      City      `gorm:"foreignKey:CityID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
