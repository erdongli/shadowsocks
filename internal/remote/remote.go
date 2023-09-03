package remote

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

type Remote struct {
	ln net.Listener
}

func New(port string) (*Remote, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, err
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

	sconn := shadow.Shadow(conn, cfg.PSK, cfg.AEADConfig)

	addr, err := socks.Address(sconn)
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

	log.Printf("start relaying between %s <-> %s", conn.RemoteAddr(), fconn.RemoteAddr())
	duplex.Relay(fconn, sconn)
}
