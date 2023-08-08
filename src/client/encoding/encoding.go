package encoding

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// Encoder
// Main operating struct. To initialize call New()
type Encoder struct {
	crypt        cipher.Block
	secret       []byte
	randomReader io.Reader
}

func New(secret string) *Encoder {
	key := sha256.Sum256([]byte(secret))
	aesBlock, _ := aes.NewCipher(key[:])
	return &Encoder{
		crypt:        aesBlock,
		secret:       key[:],
		randomReader: rand.Reader,
	}
}

// Decode
// Decodes string using given secret
func (e *Encoder) Decode(in []byte) (out string) {
	dst := make([]byte, aes.BlockSize)
	e.crypt.Decrypt(dst, in)
	return string(dst)
}

// Encode
// Encodes string using given secret
func (e *Encoder) Encode(in string) (out []byte) {
	dst := make([]byte, aes.BlockSize)
	e.crypt.Encrypt(dst, []byte(in))
	return dst
}
