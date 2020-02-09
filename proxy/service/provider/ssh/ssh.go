package ssh

import (
	"fmt"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/imsilence/kaleidoscope/proxy/config"

	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/backend"
	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/proxy"
	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/sshd"
)

type Ssh struct {
	serviceConfig config.ServiceConfig
	addr          string
	server        *sshd.Server
	running       bool
}

func (s *Ssh) Name() string {
	return "ssh"
}

func (s *Ssh) Init(addr string, serviceConfig config.ServiceConfig) {
	s.serviceConfig = serviceConfig
	s.addr = net.JoinHostPort(addr, strconv.Itoa(serviceConfig.Port))
	s.server = sshd.NewServer(s.addr)
}

func (s *Ssh) ListenAndServe() error {
	logrus.Debug("ssh server started on:", s.addr)
	s.running = true
	err := s.server.Listen()
	if err != nil {
		return err
	}
	for {
		if !s.running {
			break
		}
		conn, err := s.server.Accept()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("error accept connection")
			continue
		}
		go func(conn *sshd.SshConn) {
			defer conn.Close()
			backend, err := backend.NewBackend(s.serviceConfig, conn.Ctx)
			if err != nil {
				fmt.Println(err)
			} else {
				defer backend.Close()
				proxy.NewProxy(conn, backend).Run()
			}
		}(conn)

	}
	return nil
}

func (s *Ssh) Shutdown() error {
	s.running = false
	return s.server.Close()
}
