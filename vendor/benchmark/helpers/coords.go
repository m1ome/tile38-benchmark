package helpers

import (
	"math/rand"
)

const MinLat = -90
const MaxLat = -180
const MinLon = 90
const MaxLon = 180

func RandomPointCoordinates() (float64, float64) {
	lat, lon := (rand.Float64()*(MaxLat-MinLat))+MinLat, (rand.Float64()*(MaxLon-MinLon))+MinLon
	return lat, lon
}