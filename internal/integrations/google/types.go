package google

import "time"

type WeatherResponse struct {
	CurrentTime time.Time `json:"currentTime"`
	TimeZone    struct {
		Id string `json:"id"`
	} `json:"timeZone"`
	//IsDaytime        bool `json:"isDaytime"`
	WeatherCondition struct {
		IconBaseUri string   `json:"iconBaseUri"`
		Description struct { // ++
			Text         string `json:"text"`
			LanguageCode string `json:"languageCode"`
		} `json:"description"`
		Type string `json:"type"`
	} `json:"weatherCondition"`
	Temperature struct { // ++
		Degrees float64 `json:"degrees"`
		Unit    string  `json:"unit"`
	} `json:"temperature"`
	RelativeHumidity int `json:"relativeHumidity"` // ++
	//UvIndex          int `json:"uvIndex"`
	//Precipitation    struct {
	//	Probability struct {
	//		Percent int    `json:"percent"`
	//		Type    string `json:"type"`
	//	} `json:"probability"`
	//	Qpf struct {
	//		Quantity float64 `json:"quantity"`
	//		Unit     string  `json:"unit"`
	//	} `json:"qpf"`
	//} `json:"precipitation"`
	//ThunderstormProbability int `json:"thunderstormProbability"`
	//AirPressure             struct {
	//	MeanSeaLevelMillibars float64 `json:"meanSeaLevelMillibars"`
	//} `json:"airPressure"`
	//Wind struct {
	//	Direction struct {
	//		Degrees  int    `json:"degrees"`
	//		Cardinal string `json:"cardinal"`
	//	} `json:"direction"`
	//	Speed struct {
	//		Value int    `json:"value"`
	//		Unit  string `json:"unit"`
	//	} `json:"speed"`
	//	Gust struct {
	//		Value int    `json:"value"`
	//		Unit  string `json:"unit"`
	//	} `json:"gust"`
	//} `json:"wind"`
	//Visibility struct {
	//	Distance int    `json:"distance"`
	//	Unit     string `json:"unit"`
	//} `json:"visibility"`
	//CloudCover               int `json:"cloudCover"`
	//CurrentConditionsHistory struct {
	//	TemperatureChange struct {
	//		Degrees float64 `json:"degrees"`
	//		Unit    string  `json:"unit"`
	//	} `json:"temperatureChange"`
	//	MaxTemperature struct {
	//		Degrees float64 `json:"degrees"`
	//		Unit    string  `json:"unit"`
	//	} `json:"maxTemperature"`
	//	MinTemperature struct {
	//		Degrees float64 `json:"degrees"`
	//		Unit    string  `json:"unit"`
	//	} `json:"minTemperature"`
	//	Qpf struct {
	//		Quantity int    `json:"quantity"`
	//		Unit     string `json:"unit"`
	//	} `json:"qpf"`
	//} `json:"currentConditionsHistory"`
}

type CityInfo struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	GooglePlaceID string  `json:"googlePlaceId"`
}
