package benchmark

import (
	"net"
	"fmt"
	"io"
	"strconv"
	"bufio"
	"errors"
)

const MaxMessageSize = 0x1FFFFFFF

type Connection struct {
	connection net.Conn
	reader     *bufio.Reader
	options    ConnectionOptions
}

type ConnectionOptions struct {
	hostname  string
	port      string
	socket    string
	password  string
	keepalive bool
}

func NewConnection(opts ConnectionOptions) *Connection {
	if opts.socket == "" && opts.hostname == "" {
		opts.hostname = "127.0.0.1"
	}

	if opts.socket == "" && opts.port == "" {
		opts.port = "9851"
	}

	return &Connection{options: opts}
}

func PackMessage(w io.Writer, message []byte) error {
	h := []byte("$" + strconv.FormatUint(uint64(len(message)), 10) + " ")
	b := make([]byte, len(h) + len(message) + 2)
	copy(b, h)
	copy(b[len(h):], message)
	b[len(b) - 2] = '\r'
	b[len(b) - 1] = '\n'
	_, err := w.Write(b)
	return err
}

func UnpackMessage(w io.Writer, r *bufio.Reader) ([]byte, error) {
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

func (c *Connection) Dial() error {
	address := fmt.Sprintf("%s:%s", c.options.hostname, c.options.port)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return err
	}

	c.connection = conn
	c.reader = bufio.NewReader(conn)
	return nil
}

func (c *Connection) Close() error {
	return c.connection.Close()
}

func (c *Connection) Do(command string) ([]byte, error) {
	if err := PackMessage(c.connection, []byte(command)); err != nil {
		return nil, err
	}

	m, err := UnpackMessage(c.connection, c.reader)
	if err != nil {
		return nil, err
	}

	return m, nil
}
