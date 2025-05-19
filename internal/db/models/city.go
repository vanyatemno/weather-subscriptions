package models

import "strconv"

type City struct {
	ID            string  `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name          string  `gorm:"not null;unique"`
	Longitude     float64 `gorm:"not null"`
	Latitude      float64 `gorm:"not null"`
	GooglePlaceID string  `gorm:"not null;unique"`
	Users         []User  `gorm:"foreignKey:CityID"`
}

type Coordinates struct {
	Long string
	Lat  string
}

func (c *City) GetStringCoordinates() Coordinates {
	longitudeStr := strconv.FormatFloat(c.Longitude, 'f', -1, 64)
	latitudeStr := strconv.FormatFloat(c.Latitude, 'f', -1, 64)
	return Coordinates{
		Long: longitudeStr,
		Lat:  latitudeStr,
	}
}
