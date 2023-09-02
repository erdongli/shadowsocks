package remote

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

type Remote struct {
	ln net.Listener
}

func New(port string) (*Remote, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	return &Remote{
		ln: ln,
	}, nil
}

func (r *Remote) Serve() error {
	log.Printf("accepting connection on address %s", r.ln.Addr())

	for {
		conn, err := r.ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	addr, err := socks.Address(conn)
	if err != nil {
		log.Printf("failed to read address: %v", err)
		return
	}

	fconn, err := net.Dial(network, addr.String())
	if err != nil {
		log.Printf("failed to create forward connection: %v", err)
		return
	}
	defer fconn.Close()

	duplex.Relay(fconn, conn)
}
