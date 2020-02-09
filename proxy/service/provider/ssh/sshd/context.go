package sshd

import (
	"context"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type contextKey string

var (
	contextKeyServer        = contextKey("server")
	contextKeyUser          = contextKey("user")
	contextKeyPassword      = contextKey("password")
	contextKeyPublicKey     = contextKey("public-key")
	contextKeySessionID     = contextKey("session-id")
	contextKeyRemoteAddr    = contextKey("remote-addr")
	contextKeyLocalAddr     = contextKey("local-addr")
	contextKeyClientVersion = contextKey("client-version")
	contextKeyServerVersion = contextKey("server-version")
)

type Context interface {
	context.Context
	sync.Locker

	User() string
	Password() string
	PublicKey() string
	SessionID() string
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	ClientVersion() string
	ServerVersion() string
	SetValue(key, value interface{})
}

type sshContext struct {
	context.Context
	*sync.Mutex
}

func newContext(srv *Server) (*sshContext, context.CancelFunc) {
	ictx, cancel := context.WithCancel(context.Background())
	ctx := &sshContext{ictx, &sync.Mutex{}}

	ctx.SetValue(contextKeyServer, srv)
	return ctx, cancel
}

func applyConnMetadata(ctx Context, metadata ssh.ConnMetadata) {
	if ctx.Value(contextKeySessionID) != nil {
		return
	}
	ctx.SetValue(contextKeySessionID, string(metadata.SessionID()))
	ctx.SetValue(contextKeyUser, metadata.User())
	ctx.SetValue(contextKeyRemoteAddr, metadata.RemoteAddr())
	ctx.SetValue(contextKeyLocalAddr, metadata.LocalAddr())
	ctx.SetValue(contextKeyClientVersion, string(metadata.ClientVersion()))
	ctx.SetValue(contextKeyServerVersion, string(metadata.ServerVersion()))
}

func (ctx *sshContext) SetValue(key, value interface{}) {
	ctx.Context = context.WithValue(ctx.Context, key, value)
}

func (ctx *sshContext) User() string {
	return ctx.Value(contextKeyUser).(string)
}

func (ctx *sshContext) Password() string {
	return ctx.Value(contextKeyPassword).(string)
}

func (ctx *sshContext) PublicKey() string {
	return ctx.Value(contextKeyPublicKey).(string)
}

func (ctx *sshContext) SessionID() string {
	return ctx.Value(contextKeySessionID).(string)
}

func (ctx *sshContext) RemoteAddr() net.Addr {
	return ctx.Value(contextKeyRemoteAddr).(net.Addr)
}

func (ctx *sshContext) LocalAddr() net.Addr {
	return ctx.Value(contextKeyLocalAddr).(net.Addr)
}

func (ctx *sshContext) ClientVersion() string {
	return ctx.Value(contextKeyClientVersion).(string)
}

func (ctx *sshContext) ServerVersion() string {
	return ctx.Value(contextKeyServerVersion).(string)
}
