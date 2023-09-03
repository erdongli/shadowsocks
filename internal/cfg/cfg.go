package cfg

import "github.com/erdongli/shadowsocks-go/internal/shadow"

const (
	LocalPort = "1080"

	RemoteHost = "localhost"
	RemotePort = "1081"
)

var (
	AEADConfig = shadow.ChaCha20Poly1305
	PSK        = []byte("S4nswHpSBbzGK7a8efx10awulgSHvpQB")
)
