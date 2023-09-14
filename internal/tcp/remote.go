package tcp

import (
	"net"

	"github.com/erdongli/shadowsocks/internal/log"
	"github.com/erdongli/shadowsocks/internal/shadow"
	"github.com/erdongli/shadowsocks/internal/socks"
)

// Remote is a TCP-based ss-remote.
type Remote struct {
	ln  net.Listener
	psk []byte
	cfg shadow.AEADConfig
}

// NewRemote creates an TCP-based ss-remote.
func NewRemote(port, key string, cfg shadow.AEADConfig) (*Remote, error) {
	ln, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}

	return &Remote{
		ln:  ln,
		psk: cfg.PSK(key),
		cfg: cfg,
	}, nil
}

// Serve blocks and starts serving incoming connections.
func (r *Remote) Serve() {
	log.Info("accepting connection on address %s", r.ln.Addr())

	for {
		conn, err := r.ln.Accept()
		if err != nil {
			log.Warn("failed to accept connection: %v", err)
			continue
		}

		go r.handle(conn)
	}
}

func (r *Remote) handle(conn net.Conn) {
	defer conn.Close()

	sconn := shadow.Shadow(conn, r.psk, r.cfg)

	addr, err := socks.ReadSocksAddr(sconn)
	if err != nil {
		log.Error("failed to read address: %v", err)
		return
	}

	fconn, err := net.Dial(network, addr.String())
	if err != nil {
		log.Warn("failed to create forward connection: %v", err)
		return
	}
	defer fconn.Close()

	log.Debug("connecting to %s for %s", addr, conn.RemoteAddr())
	relay(fconn, sconn)
}
