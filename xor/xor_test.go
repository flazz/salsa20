package xor

import (
	"crypto/rand"
	"testing"
)

const K = 1024

type sliceFun func([]byte, []byte, []byte) int

func TestXor(t *testing.T) {
	a := make([]byte, K)
	if _, err := rand.Read(a); err != nil {
		t.Error(err)
	}

	b := make([]byte, K*2)
	if _, err := rand.Read(b); err != nil {
		t.Error(err)
	}

	out := make([]byte, K+100)

	table := map[string](sliceFun){
		//"simd":  xorSimd,
		"näive": xorBytes,
		"fast":  xorWords,
	}

	for name, xor := range table {
		n := xor(out, a, b)
		if n != len(a) {
			t.Errorf("%v xor length: want %v; got %v", name, len(out), n)
		}

		for i := range a {
			if out[i] != a[i]^b[i] {
				t.Errorf("%v xor not working", name)
			}
		}

	}
}

func BenchmarkNäive(b *testing.B) {
	benchmarkSliceFun(xorBytes, b)
}

func BenchmarkFast(b *testing.B) {
	benchmarkSliceFun(xorWords, b)
}

// func BenchmarkSimd(b *testing.B) {
// 	benchmarkSliceFun(xorSimd, b)
// }

func benchmarkSliceFun(fn sliceFun, b *testing.B) {
	a := make([]byte, K)
	if _, err := rand.Read(a); err != nil {
		b.Error(err)
	}
	out := make([]byte, K)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(out, a, a)
	}
}
