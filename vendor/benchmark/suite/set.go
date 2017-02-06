package suite

import (
	"benchmark/connection"
	"fmt"
	"errors"
	"benchmark/helpers"
	"math/rand"
)

const testPrefixSet = "set"

type SetCommand struct {
	T GeoType
}

func (c *SetCommand) Fire(conn connection.ConnectorReadWriter) error {
	var command string

	switch c.T {
	case Point:
		lat, lon := helpers.RandomPointCoordinates()
		command = fmt.Sprintf("SET %s %s_%d POINT %f %f", BenchmarkPrefix, testPrefixSet, rand.Intn(1000), lat, lon)
	case Geohash:
		hash, _ := helpers.RandomHashCoordinates()
		command = fmt.Sprintf("SET %s %s_%d HASH %s %d", BenchmarkPrefix, testPrefixSet, rand.Intn(1000), hash, 10)
	case Object:
		json := helpers.RandomJsonCoordinates()
		command = fmt.Sprintf("SET %s %s_%d OBJECT %s", BenchmarkPrefix, testPrefixSet, rand.Intn(1000), json)
	case Bounds:
		lat1, lon1, lat2, lon2 := helpers.RandomBoundsCoordinates()
		command = fmt.Sprintf("SET %s %s_%d POINT %f %f %f %f", BenchmarkPrefix, testPrefixSet, rand.Intn(1000), lat1, lon1, lat2, lon2)
	}

	return conn.Write(command)
}

func  (c *SetCommand) Match(conn connection.ConnectorReadWriter) error {
	data, err := conn.Read()

	if err != nil{
		return err
	}

	if len(data) == 0 {
		return errors.New("Empty response from Tile38 server")
	}

	return nil
}

func (c *SetCommand) Up(conn connection.ConnectorReadWriter) error {
	return nil
}

func (c *SetCommand) Down(conn connection.ConnectorReadWriter) error {
	return nil
}