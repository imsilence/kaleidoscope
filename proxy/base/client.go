package base

import (
	"io"
	"net"
)

type Client interface {
	io.ReadWriteCloser
	ID() string
	Addr() string
	Connection() net.Conn
}
