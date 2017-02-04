package suite

import (
	"benchmark/connection"
	"errors"
)

type GetCommand struct {}


func (c *GetCommand) Fire(conn connection.ConnectorReadWriter) error {
	return conn.Write("GET fleet truck")
}

func  (c *GetCommand) Match(conn connection.ConnectorReadWriter) error {
	data, err := conn.Read()

	if err != nil{
		return err
	}

	if len(data) == 0 {
		return errors.New("Empty response from Tile38 server")
	}

	return nil
}