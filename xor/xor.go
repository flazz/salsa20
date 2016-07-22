package xor

import (
	"runtime"
	"unsafe"
)

// largely copied from https://golang.org/src/crypto/cipher/xor.go
const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64"

// Xor performs an xor such that dst[i] = a[i] ^ b[i];
// it returns the number of bytes operated on
func Xor(dst, a, b []byte) int {
	if supportsUnaligned {
		return xorWords(a, b, dst)
	}

	return xorBytes(a, b, dst)
}

func xorBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i] // XORQ
	}
	return n
}

// TODO try out PXOR in asm
//func xorSimd(dst, a, b []byte) int

func xorWords(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] ^ bw[i]
		}
	}

	for i := (n - n%wordSize); i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}

	return n
}
