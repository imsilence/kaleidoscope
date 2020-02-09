package ssh

import (
	"fmt"
	"net"
)

type sshClient struct {
	conn net.Conn
}

func newSshClient(conn net.Conn) *sshClient {
	return &sshClient{conn}
}

func (c *sshClient) Read(p []byte) (int, error) {
	return 0, nil
}

func (c *sshClient) Write(p []byte) (int, error) {
	return 0, nil
}

func (c *sshClient) Close() error {
	return nil
}

func (c *sshClient) ID() string {
	return fmt.Sprintf("%s-%s", c.conn.LocalAddr().String(), c.conn.RemoteAddr().String())
}

func (c *sshClient) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *sshClient) Connection() net.Conn {
	return c.conn
}
