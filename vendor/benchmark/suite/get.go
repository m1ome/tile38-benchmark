package suite

import (
	"benchmark/connection"
	"errors"
	"fmt"
	"benchmark/helpers"
	"strings"
	"math/rand"
)

const maxId = 100
const testPrefixGet = "get"

type GetCommand struct {
	T GeoType
}

func (c *GetCommand) Up(conn connection.ConnectorReadWriter) error {
	var err error

	for i := 0; i < maxId; i++ {
		lat, lon := helpers.RandomPointCoordinates()
		c := fmt.Sprintf("SET %s %s_%d POINT %f %f", BenchmarkPrefix, testPrefixGet, i, lat, lon)
		err = conn.Write(c)

		if err != nil {
			return err
		}

		m, err := conn.Read()
		if err != nil {
			return err
		}

		if !strings.Contains(string(m), "\"ok\":true") {
			return errors.New(fmt.Sprintf("Error: %s", string(m)))
		}
	}

	return nil
}

func (c *GetCommand) Down(conn connection.ConnectorReadWriter) error {
	var err error

	for i := 0; i < maxId; i++ {
		c := fmt.Sprintf("DEL %s %s_%d", BenchmarkPrefix, testPrefixGet, i)
		err = conn.Write(c)

		if err != nil {
			return err
		}

		m, err := conn.Read()
		if err != nil {
			return err
		}

		if !strings.Contains(string(m), "\"ok\":true") {
			return errors.New(fmt.Sprintf("Error: %s", string(m)))
		}
	}

	return nil
}

func (c *GetCommand) Fire(conn connection.ConnectorReadWriter) error {
	command := fmt.Sprintf("GET %s %s_%d %s", BenchmarkPrefix, testPrefixGet, rand.Intn(maxId), c.T)
	return conn.Write(command)
}

func  (c *GetCommand) Match(conn connection.ConnectorReadWriter) error {
	data, err := conn.Read()

	if err != nil{
		return err
	}

	if len(data) == 0 {
		return errors.New("Empty response from Tile38 server")
	}

	if !strings.Contains(string(data), "\"ok\":true") {
		return errors.New(fmt.Sprintf("Error: %s", string(data)))
	}

	return nil
}