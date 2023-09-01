package local

import (
	"fmt"
	"log"
	"net"

	"github.com/erdongli/shadowsocks-go/internal/duplex"
	"github.com/erdongli/shadowsocks-go/internal/socks"
)

const (
	network = "tcp"
)

type Local struct {
	ln net.Listener
}

func New(p string) (*Local, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", p))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", p, err)
	}

	return &Local{
		ln: ln,
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

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	addr, port, err := socks.Handshake(conn)
	if err != nil {
		log.Printf("failed to perform handshake: %v", err)
		return
	}

	fconn, err := net.Dial(network, net.JoinHostPort(addr, port))
	if err != nil {
		log.Printf("failed to create forward tunnel: %v", err)
		return
	}
	defer fconn.Close()

	duplex.Relay(fconn, conn)
}
