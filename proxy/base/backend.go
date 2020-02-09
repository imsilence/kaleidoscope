package base

import (
	"io"
)

type Backend interface {
	io.ReadWriteCloser
}
