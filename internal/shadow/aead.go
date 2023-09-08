package shadow

import (
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"io"

	"github.com/erdongli/shadowsocks-go/internal/math"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

const (
	maxPayloadLen = 0x3FFF
	minSaltLen    = 16
)

var (
	info = []byte("ss-subkey")

	ChaCha20Poly1305 = AEADConfig{
		KeySize:  chacha20poly1305.KeySize,
		SaltSize: 32,
		TagSize:  chacha20poly1305.Overhead,
		New:      chacha20poly1305.New,
		PSK: func(k string) []byte {
			psk := sha256.Sum256([]byte(k))
			return psk[:]
		},
	}
)

type AEADConfig struct {
	KeySize, SaltSize, TagSize int

	New func(key []byte) (cipher.AEAD, error)
	PSK func(k string) []byte
}

func newAEAD(sr io.Reader, psk []byte, cfg AEADConfig) (cipher.AEAD, []byte, error) {
	s := make([]byte, math.Max(minSaltLen, cfg.KeySize))

	if _, err := io.ReadFull(sr, s); err != nil {
		return nil, nil, err
	}

	hr := hkdf.New(sha1.New, psk, s, info)
	sk := make([]byte, len(psk))

	if _, err := io.ReadFull(hr, sk); err != nil {
		return nil, nil, err
	}

	aead, err := cfg.New(sk)
	if err != nil {
		return nil, nil, err
	}

	return aead, s, nil
}
