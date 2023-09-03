package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	version = 0x5

	maxNMethods        = 255
	methodNoAuthN      = byte(0x0)
	methodNoAcceptable = byte(0xff)

	cmdConnect      = 0x1
	cmdUDPAssociate = 0x3

	rsv = 0x0

	repSucceeded       = 0x0
	repCMDNotSupported = 0x7

	atypIPv4 = 0x1
	atypFQDN = 0x3
	atypIPv6 = 0x4

	portLen     = 2
	ipv4AddrLen = 1 + net.IPv4len + portLen
	ipv6AddrLen = 1 + net.IPv6len + portLen
	maxAddrLen  = 1 + 1 + 255 + portLen
)

type Addr interface {
	Bytes() []byte
	String() string
}

// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+
//
// +----+--------+
// |VER | METHOD |
// +----+--------+
// | 1  |   1    |
// +----+--------+
//
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+

// +----+-----+-------+------+----------+----------+
// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
func Handshake(rw io.ReadWriter) (Addr, error) {
	buf := [maxNMethods]byte{}

	if _, err := io.ReadFull(rw, buf[:2]); err != nil {
		return nil, err
	}

	if buf[0] != version {
		return nil, fmt.Errorf("invalid version 0x%x", buf[0])
	}
	n := int(buf[1])

	if _, err := io.ReadFull(rw, buf[:n]); err != nil {
		return nil, err
	}

	m := methodNoAcceptable
	for _, v := range buf[:n] {
		if v == methodNoAuthN {
			m = methodNoAuthN
			break
		}
	}

	if _, err := rw.Write([]byte{version, m}); err != nil {
		return nil, err
	}

	if m == methodNoAcceptable {
		return nil, fmt.Errorf("no acceptable methods")
	}

	if _, err := io.ReadFull(rw, buf[:3]); err != nil {
		return nil, err
	}

	if buf[0] != version {
		return nil, fmt.Errorf("invalid version 0x%x", buf[0])
	}
	cmd := buf[1]

	addr, err := Address(rw)
	if err != nil {
		return nil, err
	}

	switch cmd {
	case cmdConnect:
		if _, err := rw.Write([]byte{version, repSucceeded, rsv, atypIPv4, 0, 0, 0, 0, 0, 0}); err != nil {
			return nil, err
		}
	case cmdUDPAssociate:
		fallthrough
	default:
		rw.Write([]byte{version, repCMDNotSupported, rsv})
		return nil, fmt.Errorf("unsupported command 0x%x", cmd)
	}

	return addr, nil
}

// +------+----------+------+
// | ATYP |   ADDR   | PORT |
// +------+----------+------+
// |  1   | Variable |  2   |
// +------+----------+------+
func Address(r io.Reader) (Addr, error) {
	buf := [maxAddrLen]byte{}

	if _, err := io.ReadFull(r, buf[:1]); err != nil {
		return nil, err
	}

	switch buf[0] {
	case atypIPv4:
		if _, err := io.ReadFull(r, buf[1:ipv4AddrLen]); err != nil {
			return nil, err
		}

		return ipv4Addr(buf[:ipv4AddrLen]), nil
	case atypFQDN:
		if _, err := io.ReadFull(r, buf[1:1+1]); err != nil {
			return nil, err
		}
		l := 2 + int(buf[1]) + portLen

		if _, err := io.ReadFull(r, buf[2:l]); err != nil {
			return nil, err
		}

		return fqdnAddr(buf[:l]), nil
	case atypIPv6:
		if _, err := io.ReadFull(r, buf[1:ipv6AddrLen]); err != nil {
			return nil, err
		}

		return ipv6Addr(buf[:ipv6AddrLen]), nil
	default:
		return nil, fmt.Errorf("unsupported address type: 0x%x", buf[0])
	}
}

type ipv4Addr []byte

func (a ipv4Addr) Bytes() []byte {
	return a
}

func (a ipv4Addr) String() string {
	h := net.IP(a[1 : 1+net.IPv4len]).String()
	return net.JoinHostPort(h, port(a))
}

type fqdnAddr []byte

func (a fqdnAddr) Bytes() []byte {
	return a
}

func (a fqdnAddr) String() string {
	h := string(a[2 : len(a)-portLen])
	return net.JoinHostPort(h, port(a))
}

type ipv6Addr []byte

func (a ipv6Addr) Bytes() []byte {
	return a
}

func (a ipv6Addr) String() string {
	h := net.IP(a[1 : 1+net.IPv6len]).String()
	return net.JoinHostPort(h, port(a))
}

func port(b []byte) string {
	return strconv.Itoa(int(binary.BigEndian.Uint16(b[len(b)-portLen:])))
}
