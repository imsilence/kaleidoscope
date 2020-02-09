package proxy

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/backend"
	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh/sshd"
)

func sessionHandler(conn *sshd.SshConn, newChannel ssh.NewChannel, backend *backend.Backend) {
	channel, reqs, err := newChannel.Accept()
	if err != nil {
		return
	}
	sshSession, err := backend.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	session := Session{
		conn:    conn,
		channel: channel,
		session: sshSession,
	}
	sshSession.RequestPty("kk", 1024, 1024, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	})
	session.HandleRequest(reqs)
}

type Session struct {
	conn    *sshd.SshConn
	channel ssh.Channel
	session *ssh.Session
}

func (s *Session) HandleRequest(reqs <-chan *ssh.Request) {
	for req := range reqs {
		switch req.Type {
		// case "x11-req":
		// case "pty-req":
		// case "exec":
		// case "env":
		case "shell":
			in, _ := s.session.StdinPipe()
			out, _ := s.session.StdoutPipe()

			name := strings.ReplaceAll(fmt.Sprintf("%s-%s-%s",
				strconv.FormatInt(time.Now().UnixNano(), 10),
				s.conn.Ctx.LocalAddr(),
				s.conn.Ctx.RemoteAddr(),
			), ":", "##")

			file, err := os.Create(path.Join("logs", name))
			defer file.Close()
			fmt.Println(err)
			c := io.MultiWriter(s.channel, file)
			// b := io.MultiWriter(in, os.Stdout)

			s.session.Shell()
			go func() {
				_, err := io.Copy(c, out)
				if err != nil {
					status := struct{ Status uint32 }{uint32(0)}
					s.channel.SendRequest("exit-status", false, ssh.Marshal(&status))
				}
			}()
			go func() {
				_, err := io.Copy(in, s.channel)
				if err != nil {
					status := struct{ Status uint32 }{uint32(0)}
					s.channel.SendRequest("exit-status", false, ssh.Marshal(&status))
				}
			}()
		default:
			logrus.WithFields(logrus.Fields{
				"type": req.Type,
			}).Warn("no handle to session req")
			req.Reply(false, nil)
		}
	}
}
