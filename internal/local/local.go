package local

import (
	"fmt"
	"log"
	"net"

	"github.com/erdongli/shadowsocks-go/internal/socks"
)

const (
	network = "tcp"
)

type Local struct {
	ln net.Listener
}

func New(p int) (*Local, error) {
	ln, err := net.Listen(network, fmt.Sprintf(":%d", p))
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

	log.Printf("tunneling to %s:%d", addr, port)
}
