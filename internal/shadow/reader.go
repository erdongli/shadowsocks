package shadow

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/erdongli/shadowsocks/internal/math"
)

type reader struct {
	r      io.Reader
	aead   cipher.AEAD
	nonce  []byte
	buf    []byte
	offset int
}

func newReader(r io.Reader, psk []byte, cfg AEADConfig) (io.Reader, error) {
	aead, _, err := newAEAD(r, psk, cfg)
	if err != nil {
		return nil, err
	}

	return &reader{
		r:      r,
		aead:   aead,
		nonce:  make([]byte, aead.NonceSize()),
		buf:    []byte{},
		offset: 0,
	}, nil
}

func (r *reader) Read(b []byte) (int, error) {
	if r.offset == len(r.buf) {
		buf := make([]byte, maxPayloadSize+r.aead.Overhead())

		if _, err := io.ReadFull(r.r, buf[:2+r.aead.Overhead()]); err != nil {
			return 0, err
		}

		if _, err := r.aead.Open(buf[:0], r.nonce, buf[:2+r.aead.Overhead()], nil); err != nil {
			return 0, err
		}
		math.IncrLittleEndian(r.nonce)

		l := int(binary.BigEndian.Uint16(buf))
		if l > maxPayloadSize {
			return 0, fmt.Errorf("invalid payload size %d", l)
		}

		if _, err := io.ReadFull(r.r, buf[:l+r.aead.Overhead()]); err != nil {
			return 0, err
		}

		r.buf, r.offset = make([]byte, l), 0
		if _, err := r.aead.Open(r.buf[:0], r.nonce, buf[:l+r.aead.Overhead()], nil); err != nil {
			return 0, err
		}
		math.IncrLittleEndian(r.nonce)
	}

	n := copy(b, r.buf[r.offset:])
	r.offset += n

	return n, nil
}
