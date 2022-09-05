package grpcquic

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

type Conn struct {
	conn   quic.Connection
	stream quic.Stream
}

func NewConn(conn quic.Connection) (net.Conn, error) {
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}

	return &Conn{conn, stream}, nil
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.stream.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return c.stream.Write(b)
}

func (c *Conn) Close() error {
	c.stream.Close()
	return c.conn.CloseWithError(0, "")
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)

}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}

type Listener struct {
	ql quic.Listener
}

func Listen(ql quic.Listener) net.Listener {
	return &Listener{ql}
}

func (l *Listener) Accept() (net.Conn, error) {
	conn, err := l.ql.Accept(context.Background())
	if err != nil {
		return nil, err
	}

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}

	return &Conn{conn, stream}, nil
}

func (l *Listener) Close() error {
	return l.ql.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.ql.Addr()
}

func NewQuicDialer(conf *tls.Config) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, target string) (net.Conn, error) {
		conn, err := quic.DialAddrContext(ctx, target, conf, &quic.Config{})
		if err != nil {
			return nil, err
		}

		return NewConn(conn)
	}
}
