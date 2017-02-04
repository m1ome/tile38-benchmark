package suite

import (
	"benchmark/connection"
	"fmt"
	"errors"
)

type SetCommand struct {}

func (c *SetCommand) Fire(conn connection.ConnectorReadWriter) error {
	lat, lon := randomPointCoordinates()
	command := fmt.Sprintf("SET fleet truck POINT %f %f", lat, lon)

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