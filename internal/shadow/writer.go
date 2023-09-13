package shadow

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"io"

	"github.com/erdongli/shadowsocks/internal/math"
)

type writer struct {
	w     io.Writer
	aead  cipher.AEAD
	nonce []byte
}

func newWriter(w io.Writer, psk []byte, cfg AEADConfig) (io.Writer, error) {
	aead, salt, err := newAEAD(rand.Reader, psk, cfg)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(salt); err != nil {
		return nil, err
	}

	return &writer{
		w:     w,
		aead:  aead,
		nonce: make([]byte, aead.NonceSize()),
	}, nil
}

func (w *writer) Write(b []byte) (int, error) {
	n, buf := 0, make([]byte, maxPayloadSize+w.aead.Overhead())

	for n < len(b) {
		l := math.Min(maxPayloadSize, len(b)-n)
		binary.BigEndian.PutUint16(buf, uint16(l))

		w.aead.Seal(buf[:0], w.nonce, buf[:2], nil)
		math.IncrLittleEndian(w.nonce)

		if _, err := w.w.Write(buf[:2+w.aead.Overhead()]); err != nil {
			return n, err
		}

		w.aead.Seal(buf[:0], w.nonce, b[n:n+l], nil)
		math.IncrLittleEndian(w.nonce)

		if _, err := w.w.Write(buf[:l+w.aead.Overhead()]); err != nil {
			return n, err
		}

		n += l
	}

	return n, nil
}
