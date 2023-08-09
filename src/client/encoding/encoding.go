package encoding

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"unicode/utf8"
)

var ErrWrongKey = errors.New("couldn't be decoded with given secret")

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
func (e *Encoder) Decode(in []byte) (out string, err error) {
	dst := make([]byte, aes.BlockSize)
	e.crypt.Decrypt(dst, in)
	result := bytes.Trim(dst, "\x00")
	if !utf8.Valid(result) {
		return "", ErrWrongKey
	}
	return string(result), nil
}

// Encode
// Encodes string using given secret
func (e *Encoder) Encode(in string) (out []byte) {
	src := make([]byte, int(len(in)/aes.BlockSize+1)*aes.BlockSize)
	copy(src, in)
	dst := make([]byte, aes.BlockSize)
	e.crypt.Encrypt(dst, src)
	return dst
}
