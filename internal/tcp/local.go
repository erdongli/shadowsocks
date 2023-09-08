package tcp

import (
	"log"
	"net"

	"github.com/erdongli/shadowsocks-go/internal/shadow"
	"github.com/erdongli/shadowsocks-go/internal/socks"
)

type Local struct {
	ln    net.Listener
	rAddr string
	psk   []byte
	cfg   shadow.AEADConfig
}

func NewLocal(port, rHost, rPort string, psk []byte, cfg shadow.AEADConfig) (*Local, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}

	return &Local{
		ln:    ln,
		rAddr: net.JoinHostPort(rHost, rPort),
		psk:   psk,
		cfg:   cfg,
	}, nil
}

func (l *Local) Serve() {
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

	sconn := shadow.Shadow(fconn, l.psk, l.cfg)
	if _, err := sconn.Write(addr.Bytes()); err != nil {
		log.Printf("failed to forward destination address: %v", err)
		return
	}

	log.Printf("connecting to %s for target %s", l.rAddr, addr)
	relay(sconn, conn)
}
