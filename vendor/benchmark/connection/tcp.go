package connection

import (
	"errors"
	"net"
	"bufio"
	"io"
	"strconv"
)

const MaxMessageSize = 0x1FFFFFFF

type TCPConnection struct {
	url string
}

type TCPConnectionInstance struct {
	conn net.Conn
	read *bufio.Reader
}

func NewTCPConnection(url string) *TCPConnection {
	return &TCPConnection{url}
}

func writeMessage(w io.Writer, message []byte) error {
	h := []byte("$" + strconv.FormatUint(uint64(len(message)), 10) + " ")
	b := make([]byte, len(h) + len(message) + 2)
	copy(b, h)
	copy(b[len(h):], message)
	b[len(b) - 2] = '\r'
	b[len(b) - 1] = '\n'
	_, err := w.Write(b)
	return err
}

func readMessage(r *bufio.Reader) ([]byte, error) {
	b,  err := r.ReadBytes(' ')

	if err != nil {
		return nil, err
	}

	if len(b) > 0 && b[0] != '$' {
		return nil, errors.New("Not a valid message")
	}

	n, err := strconv.ParseUint(string(b[1:len(b)-1]), 10, 32)
	if err != nil {
		return nil, errors.New("Invalid/Unparsable size")
	}

	if n > MaxMessageSize {
		return  nil, errors.New("Message is too long")
	}

	b = make([]byte, int(n) + 2)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}

	if b[len(b)-2] != '\r' || b[len(b)-1] != '\n' {
		return nil, errors.New("Expecting crlf suffix")
	}

	return b[:len(b)-2], nil
}

func (c *TCPConnection) Dial() (ConnectorReadWriter, error) {
	conn, err := net.Dial("tcp", c.url)
	if err != nil {
		return &TCPConnectionInstance{}, err
	}

	return &TCPConnectionInstance{conn, bufio.NewReader(conn)}, nil
}

func (c *TCPConnectionInstance) Close() error {
	return c.conn.Close()
}

func (c *TCPConnectionInstance) Write(command string) error {
	return writeMessage(c.conn, []byte(command))
}

func (c *TCPConnectionInstance) Read() ([]byte, error) {
	return readMessage(c.read)
}
