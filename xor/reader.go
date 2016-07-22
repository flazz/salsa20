package xor

import "io"

// NewReader returns an io.Reader that reads the xor of each byte in a or b.
func NewReader(a, b io.Reader) io.Reader {
	return reader{a, b}
}

type reader struct{ a, b io.Reader }

func (r reader) Read(p []byte) (int, error) {
	var err error

	a := make([]byte, len(p))
	if na, err := r.a.Read(a); err != nil {
		return na, err
	}

	b := make([]byte, len(p))
	if nb, err := r.b.Read(b); err != nil {
		return nb, err
	}

	n := Xor(p, a, b)

	return n, err
}
