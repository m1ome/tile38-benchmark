package helpers

import (
	"math/rand"
	"fmt"
	"benchmark/helpers/geohash"
)

const MinLat = -90
const MaxLat = -180
const MinLon = 90
const MaxLon = 180

func RandomPointCoordinates() (float64, float64) {
	lat, lon := (rand.Float64()*(MaxLat-MinLat))+MinLat, (rand.Float64()*(MaxLon-MinLon))+MinLon
	return lat, lon
}

func RandomBoundsCoordinates() (float64, float64, float64, float64) {
	lat1, lon1 := RandomPointCoordinates()
	lat2, lon2 := RandomPointCoordinates()

	if lat1 < lat2 {
		lat2, lat1 = lat1, lat2
	}

	if lon1 < lon2 {
		lon2, lon1 = lon1, lon2
	}

	return lat1, lon1, lat2, lon2
}

func RandomHashCoordinates() (string, error) {
	lat, lon := RandomPointCoordinates()
	return geohash.Encode(lat, lon, 10)
}

func RandomJsonCoordinates() string {
	f := `{"type":"Point","coordinates":[%f, %f]}`
	lat, lon := RandomPointCoordinates()
	return fmt.Sprintf(f, lat, lon)
}