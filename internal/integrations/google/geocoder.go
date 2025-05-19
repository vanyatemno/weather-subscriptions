package google

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
)

func fetchCityInfo(ctx context.Context, client *maps.Client, cityName string) (*CityInfo, error) {
	responses, err := client.Geocode(ctx, &maps.GeocodingRequest{Address: cityName})
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no results found for city with name: '%s'", cityName)
	}

	return &CityInfo{
		GooglePlaceID: responses[0].PlaceID,
		Latitude:      responses[0].Geometry.Location.Lat,
		Longitude:     responses[0].Geometry.Location.Lng,
	}, nil
}
