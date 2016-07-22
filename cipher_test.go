package salsa20

import (
	"reflect"
	"testing"
)

var cip16 Cipher

func init() {
	k32 := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		201, 202, 203, 204, 205, 206, 207, 208,
		209, 210, 211, 212, 213, 214, 215, 216,
	}
	nonce := [8]byte{
		101, 102, 103, 104, 105, 106, 107, 108,
	}

	var err error
	cip16, err = NewCipher(k32[:16], nonce)
	if err != nil {
		panic(err)
	}
}

func TestReadNothing(t *testing.T) {
	r := cip16.NewReader().(*reader)
	var p []byte
	n, err := r.Read(p)
	if err != nil {
		t.Error(err)
	}

	if n != 0 {
		t.Error("bytes returned: ", n)
	}

}

func TestReadLessThanBlockSize(t *testing.T) {
	r := cip16.NewReader().(*reader)
	p := make([]byte, blockSize/2)
	n, err := r.Read(p)
	if err != nil {
		t.Error(err)
	}

	if n != len(p) {
		t.Error("bytes returned: ", n)
	}

	if blockSize-len(r.block) != len(p) {
		t.Error("bytes read: ", len(p), len(r.block))
	}

	b0 := cip16.Block(0)
	if !reflect.DeepEqual(b0[:len(p)], p) {
		t.Error("oops", b0[:len(p)], p)
	}
}

func TestReadBlockSize(t *testing.T) {
	r := cip16.NewReader().(*reader)
	p := make([]byte, blockSize)
	n, err := r.Read(p)
	if err != nil {
		t.Error(err)
	}

	if n != len(p) {
		t.Error("bytes returned: ", n)
	}

	if blockSize-len(r.block) != len(p) {
		t.Error("bytes read: ", len(p), len(r.block))
	}

	b0 := cip16.Block(0)
	if !reflect.DeepEqual(b0[:len(p)], p) {
		t.Error("oops", b0[:len(p)], p)
	}
}

func TestReadMoreThanBlockSize(t *testing.T) {
	r := cip16.NewReader().(*reader)
	p := make([]byte, blockSize+10)
	n, err := r.Read(p)
	if err != nil {
		t.Error(err)
	}

	if n != len(p) {
		t.Error("bytes returned: ", n)
	}

	b0 := cip16.Block(0)
	b1 := cip16.Block(1)
	var bwant []byte
	bwant = append(bwant, b0[:]...)
	bwant = append(bwant, b1[:10]...)

	t.Log("want", len(bwant))
	t.Log("got", len(p))

	for i := range bwant {
		if bwant[i] != p[i] {
			t.Errorf("want[%d]  %v; got %v", i, bwant[i], p[i])
		}
	}

}

func TestReadMulti(t *testing.T) {
	r := cip16.NewReader().(*reader)
	p := make([]byte, 47) // prime number

	// 6400 bytes
	var bwant []byte
	for n := uint64(0); n < 100; n++ {
		b := cip16.Block(n)
		bwant = append(bwant, b[:]...)
	}
	t.Log("want", len(bwant))

	var bgot []byte
	for len(bgot) <= len(bwant) {
		_, err := r.Read(p)
		if err != nil {
			t.Error(err)
		}
		bgot = append(bgot, p...)
	}
	t.Log("got", len(bgot))

	for i := range bwant {
		if bwant[i] != bgot[i] {
			t.Errorf("want[%d]  %v; got %v", i, bwant[i], bgot[i])
		}
	}

}

func TestReadEOF(t *testing.T) {
	t.Skip("not implemented")
}

/*

func BenchmarkSliceFun(fn sliceFun, b *testing.B) {
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

*/
