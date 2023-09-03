package local

import (
	"log"
	"net"

	"github.com/erdongli/shadowsocks-go/internal/cfg"
	"github.com/erdongli/shadowsocks-go/internal/duplex"
	"github.com/erdongli/shadowsocks-go/internal/shadow"
	"github.com/erdongli/shadowsocks-go/internal/socks"
)

const (
	network = "tcp"
)

type Local struct {
	ln    net.Listener
	rAddr string
}

func New(port, rHost, rPort string) (*Local, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}

	return &Local{
		ln:    ln,
		rAddr: net.JoinHostPort(rHost, rPort),
	}, nil
}

func (l *Local) Serve() error {
	log.Printf("accepting connection on address %s", l.ln.Addr())

	for {
		conn, err := l.ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go l.handle(conn)
	}
}

func (l *Local) handle(conn net.Conn) {
	defer conn.Close()

	addr, err := socks.Handshake(conn)
	if err != nil {
		log.Printf("failed to perform handshake: %v", err)
		return
	}

	fconn, err := net.Dial(network, l.rAddr)
	if err != nil {
		log.Printf("failed to create forward connection: %v", err)
		return
	}
	defer fconn.Close()

	sconn := shadow.Shadow(fconn, cfg.PSK, cfg.AEADConfig)
	if _, err := sconn.Write(addr.Bytes()); err != nil {
		log.Printf("failed to forward destination address: %v", err)
		return
	}

	log.Printf("start relaying between %s <-> %s", conn.RemoteAddr(), fconn.RemoteAddr())
	duplex.Relay(sconn, conn)
}
