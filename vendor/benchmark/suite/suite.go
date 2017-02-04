package suite

import (
	"math/rand"
	"benchmark/connection"
)

type Runner interface {
	Fire(conn connection.ConnectorReadWriter) error
	Match(conn connection.ConnectorReadWriter) error
}

const MinLat = 0
const MaxLat = 0
const MinLon = 0
const MaxLon = 0

func randomPointCoordinates() (float64, float64) {
	lat, lon := (rand.Float64()*(MaxLat-MinLat))+MinLat, (rand.Float64()*(MaxLon-MinLon))+MinLon
	return lat, lon
}