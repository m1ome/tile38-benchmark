package suite

import (
	"benchmark/connection"
)

type Runner interface {
	Fire(conn connection.ConnectorReadWriter) error
	Match(conn connection.ConnectorReadWriter) error
	Up(conn connection.ConnectorReadWriter) error
	Down(conn connection.ConnectorReadWriter) error
}

const BenchmarkPrefix = "benchmark"

type GeoType int
const (
	Bounds GeoType = iota
	Point
	Geohash
	Object
)

func (g GeoType) String() string {
	switch g {
	case Bounds:
		return "BOUNDS"
	case Point:
		return "POINT"
	case Geohash:
		return "HASH 10"
	case Object:
		return "OBJECT"
	}

	return ""
}