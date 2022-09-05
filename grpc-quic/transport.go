package grpcquic

import (
	"context"
	"crypto/tls"
	"net"

	"google.golang.org/grpc/credentials"
)

type Info struct {
	conn *Conn
}

func NewInfo(c *Conn) *Info {
	return &Info{c}
}

func (i *Info) AuthType() string {
	return "quic-tls"
}

func (i *Info) Conn() net.Conn {
	return i.conn
}

var _ credentials.TransportCredentials = (*Credentials)(nil)

type Credentials struct {
	config           *tls.Config
	isQuicConnection bool
	serverName       string

	creds credentials.TransportCredentials
}

func NewCredentials(config *tls.Config) credentials.TransportCredentials {
	creds := credentials.NewTLS(config)
	return &Credentials{
		creds:  creds,
		config: config,
	}
}

func (pt *Credentials) ClientHandshake(ctx context.Context, authority string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if c, ok := conn.(*Conn); ok {
		pt.isQuicConnection = true
		return conn, NewInfo(c), nil
	}

	return pt.creds.ClientHandshake(ctx, authority, conn)
}

func (pt *Credentials) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if c, ok := conn.(*Conn); ok {
		pt.isQuicConnection = true
		ainfo := NewInfo(c)
		return conn, ainfo, nil
	}

	return pt.creds.ServerHandshake(conn)
}

func (pt *Credentials) Info() credentials.ProtocolInfo {
	if pt.isQuicConnection {
		return credentials.ProtocolInfo{
			ProtocolVersion:  "/quic/1.0.0",
			SecurityProtocol: "quic-tls",
			ServerName:       pt.serverName,
		}
	}

	return pt.creds.Info()
}

func (pt *Credentials) Clone() credentials.TransportCredentials {
	return &Credentials{
		config: pt.config.Clone(),
		creds:  pt.creds.Clone(),
	}
}

func (pt *Credentials) OverrideServerName(name string) error {
	pt.serverName = name
	return pt.creds.OverrideServerName(name)
}
