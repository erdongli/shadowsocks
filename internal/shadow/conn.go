package shadow

import (
	"io"
	"net"
)

// Conn is a Shadowsocks connection.
type Conn struct {
	net.Conn

	psk []byte
	cfg AEADConfig

	r io.Reader
	w io.Writer
}

// Shadow creates a Shadowsocks connection from the provided connection,
// pre-shared key, and AEAD configuration.
func Shadow(conn net.Conn, psk []byte, cfg AEADConfig) net.Conn {
	return &Conn{
		Conn: conn,
		psk:  psk,
		cfg:  cfg,
	}
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (int, error) {
	if c.r == nil {
		r, err := newReader(c.Conn, c.psk, c.cfg)
		if err != nil {
			return 0, err
		}

		c.r = r
	}

	return c.r.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (int, error) {
	if c.w == nil {
		w, err := newWriter(c.Conn, c.psk, c.cfg)
		if err != nil {
			return 0, err
		}

		c.w = w
	}

	return c.w.Write(b)
}
