package backend

import (
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/imsilence/kaleidoscope/proxy/config"
	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/sshd"
)

type Backend struct {
	*ssh.Client
}

func NewBackend(serviceConfig config.ServiceConfig, ctx sshd.Context) (*Backend, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            ctx.User(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(ctx.Password())}

	addr := net.JoinHostPort(serviceConfig.Backend.Host, strconv.Itoa(serviceConfig.Backend.Port))
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return &Backend{c}, nil
}

func (backend *Backend) Close() error {
	return nil
}
