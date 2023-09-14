package tcp

import (
	"net"

	"github.com/erdongli/shadowsocks/internal/log"
	"github.com/erdongli/shadowsocks/internal/shadow"
	"github.com/erdongli/shadowsocks/internal/socks"
)

// Local is a TCP-based ss-local.
type Local struct {
	ln    net.Listener
	raddr string
	psk   []byte
	cfg   shadow.AEADConfig
}

// NewLocal creates an TCP-based ss-local.
func NewLocal(port, raddr, key string, cfg shadow.AEADConfig) (*Local, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}

	return &Local{
		ln:    ln,
		raddr: raddr,
		psk:   cfg.PSK(key),
		cfg:   cfg,
	}, nil
}

// Serve blocks and starts serving incoming connections.
func (l *Local) Serve() {
	log.Info("accepting connection on address %s", l.ln.Addr())

	for {
		conn, err := l.ln.Accept()
		if err != nil {
			log.Warn("failed to accept connection: %v", err)
			continue
		}

		go l.handle(conn)
	}
}

func (l *Local) handle(conn net.Conn) {
	defer conn.Close()

	addr, err := socks.Handshake(conn)
	if err != nil {
		log.Warn("failed to perform handshake: %v", err)
		return
	}

	fconn, err := net.Dial(network, l.raddr)
	if err != nil {
		log.Error("failed to create forward connection: %v", err)
		return
	}
	defer fconn.Close()

	sconn := shadow.Shadow(fconn, l.psk, l.cfg)
	if _, err := sconn.Write(addr.Bytes()); err != nil {
		log.Error("failed to forward destination address: %v", err)
		return
	}

	log.Debug("connecting to %s for %s", l.raddr, addr)
	relay(sconn, conn)
}
