package models

type User struct {
	ID     string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Email  string `gorm:"not null;unique"`
	CityID string `gorm:"not null;foreignKey:CityID"`
	City   City   `gorm:"foreignKey:CityID"`
}
