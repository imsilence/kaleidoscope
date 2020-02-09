package sshd

import (
	"context"

	"golang.org/x/crypto/ssh"
)

type SshConn struct {
	ServerConn *ssh.ServerConn
	NewChannel <-chan ssh.NewChannel
	Request    <-chan *ssh.Request
	Ctx        Context
	Cancel     context.CancelFunc
}

func newConn(sc *ssh.ServerConn, nc <-chan ssh.NewChannel, r <-chan *ssh.Request, ctx Context, cancel context.CancelFunc) *SshConn {
	return &SshConn{sc, nc, r, ctx, cancel}
}

func (c *SshConn) Close() error {
	defer c.Cancel()
	return c.ServerConn.Close()
}
