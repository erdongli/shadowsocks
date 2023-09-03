package shadow

import (
	"io"
	"net"
)

type Conn struct {
	net.Conn

	psk []byte
	cfg AEADConfig

	r io.Reader
	w io.Writer
}

func Shadow(conn net.Conn, psk []byte, cfg AEADConfig) net.Conn {
	return &Conn{
		Conn: conn,
		psk:  psk,
		cfg:  cfg,
	}
}

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
