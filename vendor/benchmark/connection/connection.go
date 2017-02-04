package connection

type Connector interface {
	Dial() (ConnectorReadWriter, error)
}

type ConnectorReadWriter interface {
	Close() error
	Write(command string) error
	Read() ([]byte, error)
}