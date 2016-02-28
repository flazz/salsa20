package salsa20

import "testing"

func TestLRot(t *testing.T) {
	var table = []struct{ in, c, out word }{
		{0x0, 1, 0x0},
		{0x1, 1, 0x2},
		{0xf00abcde, 4, 0xabcdef},
	}
	for _, row := range table {
		r := lrot(row.in, row.c)
		if row.out != r {
			t.Errorf("want: %x; got %x\n", row.out, r)
		}
	}
}
