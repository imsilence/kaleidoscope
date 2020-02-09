package sshd

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type Server struct {
	addr     string
	listener net.Listener
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (srv *Server) Listen() (err error) {
	srv.listener, err = net.Listen("tcp", srv.addr)
	return
}

func (srv *Server) Accept() (*SshConn, error) {
	conn, err := srv.listener.Accept()
	if err != nil {
		return nil, err
	}
	ctx, cancel := newContext(srv)

	serverConn, chans, reqs, err := ssh.NewServerConn(conn, srv.config(ctx))
	if err != nil {
		return nil, err
	}

	applyConnMetadata(ctx, serverConn)

	return newConn(serverConn, chans, reqs, ctx, cancel), nil
}

func (srv *Server) config(ctx Context) *ssh.ServerConfig {
	conf := &ssh.ServerConfig{}
	if signer, err := generateSigner(); err == nil {
		conf.AddHostKey(signer)
	}

	conf.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		ctx.SetValue(contextKeyPassword, string(password))
		// if string(password) == "881019" {
		return &ssh.Permissions{}, nil
		// }
		// return &ssh.Permissions{}, fmt.Errorf("permission denied")
	}

	conf.PublicKeyCallback = func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		fmt.Println("public key callback")
		return &ssh.Permissions{}, fmt.Errorf("permission denied")
	}

	return conf
}

func (srv *Server) Close() error {
	srv.listener.Close()
	return nil
}
