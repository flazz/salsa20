package salsa20

import (
	"io"
	"strconv"
)

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

// KeySizeError represents a condition where the size of a given key is wrong
type KeySizeError int

func (k KeySizeError) Error() string {
	return "salsa20: invalid key size " + strconv.Itoa(int(k))
}

// Cipher represents a salsa20 stream
type Cipher struct {
	k []byte  // key
	v [8]byte // nonce
}

// TODO Cipher16 & Cipher32

// NewCipher returns a Cypher based on the key and nonce. KeySizeError is returned if the len(k) is neither 16 or 32
func NewCipher(k []byte, v [8]byte) (Cipher, error) {
	switch len(k) {
	case 16, 32:
		return Cipher{k: k, v: v}, nil

	default:
		return Cipher{}, KeySizeError(len(k))
	}
}

const (
	blockSize = 64
)

// Block returns the nth block in the cypher stream
func (c Cipher) Block(n uint64) [64]byte {
	// TODO sort out 2 << 70 limitation? maybe just use int and limit it to 2<<63?
	var nonce [16]byte

	copy(nonce[0:8], c.v[:])

	for i := uint(0); i < 8; i++ {
		nonce[i+8] = byte(n >> (i * 8))
	}

	var exp [64]byte
	switch len(c.k) {

	case 16:
		var k [16]byte
		copy(k[:], c.k)
		exp = expansion16(k, nonce)

	case 32:
		var k [32]byte
		copy(k[:], c.k)
		exp = expansion32(k, nonce)

	default:
		panic(len(c.k))
	}

	h := hash(exp)

	return h
}

// NewReader returns a new io.Reader that reads the stream
func (c Cipher) NewReader() io.Reader {
	return &reader{c, 0, nil}
}

type reader struct {
	c         Cipher
	nextBlock uint64 // TODO be able to read all 2<<70 blocks
	block     []byte
}

// Read implements io.Reader interface. Number of bytes is len(p), and error is nil
func (r *reader) Read(p []byte) (int, error) {
	// need to read nothing
	if len(p) == 0 {
		return 0, nil
	}

	var n int // bytes consumbed this read
	for n < len(p) {

		// get new block if needed
		if len(r.block) == 0 {
			b := r.c.Block(r.nextBlock)
			r.block = b[:]
			r.nextBlock++
		}

		// consume from current block
		nb := copy(p[n:], r.block)
		r.block = r.block[nb:]

		n += nb
	}

	return n, nil
}
