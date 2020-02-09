package proxy

import (
	"fmt"
	"sync"

	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/backend"
	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/sshd"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type channelHandler func(*sshd.SshConn, ssh.NewChannel, *backend.Backend)

type Proxy struct {
	conn    *sshd.SshConn
	backend *backend.Backend

	channelHandlers map[string]channelHandler
}

func NewProxy(conn *sshd.SshConn, backend *backend.Backend) *Proxy {
	channelHandlers := make(map[string]channelHandler)
	channelHandlers["session"] = sessionHandler
	return &Proxy{conn, backend, channelHandlers}
}

func (proxy *Proxy) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go proxy.handleRequest(wg)
	go proxy.handleChannel(wg)

	wg.Wait()
}

func (proxy *Proxy) handleRequest(wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range proxy.conn.Request {
		fmt.Println("req:", req.Type)
		req.Reply(false, nil)
	}
}

func (proxy *Proxy) handleChannel(wg *sync.WaitGroup) {
	defer wg.Done()
	for ch := range proxy.conn.NewChannel {
		if handler, ok := proxy.channelHandlers[ch.ChannelType()]; ok {
			go handler(proxy.conn, ch, proxy.backend)
		} else {
			logrus.WithFields(logrus.Fields{
				"channel": ch.ChannelType(),
			}).Warn("no handle to channel")
			ch.Reject(ssh.UnknownChannelType, "unsupported channel type")
		}

	}
}
