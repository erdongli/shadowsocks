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

	atypIPv4 = 0x1
	atypFQDN = 0x3
	atypIPv6 = 0x4

	methodNoAuthN      = byte(0x0)
	methodNoAcceptable = byte(0xff)

	cmdConnect      = 0x1
	cmdBind         = 0x2
	cmdUDPAssociate = 0x3

	repSucceeded       = 0x0
	repCMDNotSupported = 0x7

	rsv = 0x0

	maxBufLen = 255
	portLen   = 2
)

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
func Handshake(rw io.ReadWriter) (string, string, error) {
	buf := [maxBufLen]byte{}

	if _, err := io.ReadFull(rw, buf[:2]); err != nil {
		return "", "", fmt.Errorf("failed to read version and number of methods: %w", err)
	}
	if buf[0] != version {
		return "", "", fmt.Errorf("invalid version identifier 0x%x", buf[0])
	}
	n := int(buf[1])

	if _, err := io.ReadFull(rw, buf[:n]); err != nil {
		return "", "", fmt.Errorf("failed to read methods: %w", err)
	}
	m := methodNoAcceptable
	for _, v := range buf[:n] {
		if v == methodNoAuthN {
			m = methodNoAuthN
			break
		}
	}

	if _, err := rw.Write([]byte{version, m}); err != nil {
		return "", "", fmt.Errorf("failed to select method: %w", err)
	}

	if _, err := io.ReadFull(rw, buf[:3]); err != nil {
		return "", "", fmt.Errorf("failed to read version, command, and reserved: %w", err)
	}
	cmd := buf[1]

	addr, err := Address(rw)
	if err != nil {
		return "", "", fmt.Errorf("failed to read destination address: %w", err)
	}

	port, err := Port(rw)
	if err != nil {
		return "", "", fmt.Errorf("failed to read destination port: %w", err)
	}

	switch cmd {
	case cmdConnect:
		if _, err := rw.Write([]byte{version, repSucceeded, rsv, atypIPv4, 0, 0, 0, 0, 0, 0}); err != nil {
			return "", "", fmt.Errorf("failed to reply: %w", err)
		}
	default:
		rw.Write([]byte{version, repCMDNotSupported, rsv})
		return "", "", fmt.Errorf("unsupported command 0x%x", cmd)
	}

	return addr, port, nil
}

// +------+----------+
// | ATYP |   ADDR   |
// +------+----------+
// |  1   | Variable |
// +------+----------+
func Address(r io.Reader) (string, error) {
	buf := [maxBufLen]byte{}

	if _, err := io.ReadFull(r, buf[:1]); err != nil {
		return "", fmt.Errorf("failed to read address type: %w", err)
	}

	switch buf[0] {
	case atypIPv4:
		if _, err := io.ReadFull(r, buf[:net.IPv4len]); err != nil {
			return "", fmt.Errorf("failed to read IPv4 address: %w", err)
		}

		return net.IP(buf[:net.IPv4len]).String(), nil
	case atypFQDN:
		if _, err := io.ReadFull(r, buf[:1]); err != nil {
			return "", fmt.Errorf("failed to read length of FQDN: %w", err)
		}
		l := int(buf[0])

		if _, err := io.ReadFull(r, buf[:l]); err != nil {
			return "", fmt.Errorf("failed to read FQDN: %w", err)
		}

		return string(buf[:l]), nil
	case atypIPv6:
		if _, err := io.ReadFull(r, buf[:net.IPv6len]); err != nil {
			return "", fmt.Errorf("failed to read IPv6 address: %w", err)
		}

		return net.IP(buf[:net.IPv6len]).String(), nil
	default:
		return "", fmt.Errorf("unsupported address type: 0x%x", buf[0])
	}
}

// +------+
// | PORT |
// +------+
// |  2   |
// +------+
func Port(r io.Reader) (string, error) {
	var buf [portLen]byte

	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return "", err
	}

	return strconv.Itoa(int(binary.BigEndian.Uint16(buf[:]))), nil
}
