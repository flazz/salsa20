package salsa20

import "testing"

func TestLittleEndian(t *testing.T) {
	var table = []struct {
		in  [4]byte
		out word
	}{
		{[4]byte{0, 0, 0, 0}, 0x00000000},
		{[4]byte{86, 75, 30, 9}, 0x091e4b56},
		{[4]byte{255, 255, 255, 250}, 0xfaffffff},
	}

	for _, row := range table {
		got := littleEndian(row.in)

		if row.out != got {
			t.Errorf("want: %x; got: %x", row.out, got)
		}
	}
}

func TestLittleEndianInv(t *testing.T) {
	var table = []struct {
		in  [4]byte
		out word
	}{
		{[4]byte{0, 0, 0, 0}, 0x00000000},
		{[4]byte{86, 75, 30, 9}, 0x091e4b56},
		{[4]byte{255, 255, 255, 250}, 0xfaffffff},
	}

	for _, row := range table {
		got := littleEndianInv(row.out)

		if row.in != got {
			t.Errorf("want: %x; got: %x", row.out, got)
		}
	}
}
