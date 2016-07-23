package salsa20

import "testing"

func TestQuarterRound(t *testing.T) {
	var table = []struct{ in, out [4]word }{

		{
			[4]word{0x00000000, 0x00000000, 0x00000000, 0x00000000},
			[4]word{0x00000000, 0x00000000, 0x00000000, 0x00000000},
		},
		{
			[4]word{0x00000001, 0x00000000, 0x00000000, 0x00000000},
			[4]word{0x08008145, 0x00000080, 0x00010200, 0x20500000},
		},
		{
			[4]word{0x00000000, 0x00000001, 0x00000000, 0x00000000},
			[4]word{0x88000100, 0x00000001, 0x00000200, 0x00402000},
		},
		{
			[4]word{0x00000000, 0x00000000, 0x00000001, 0x00000000},
			[4]word{0x80040000, 0x00000000, 0x00000001, 0x00002000},
		},
		{
			[4]word{0x00000000, 0x00000000, 0x00000000, 0x00000001},
			[4]word{0x00048044, 0x00000080, 0x00010000, 0x20100001},
		},
		{
			[4]word{0xe7e8c006, 0xc4f9417d, 0x6479b4b2, 0x68c67137},
			[4]word{0xe876d72b, 0x9361dfd5, 0xf1460244, 0x948541a3},
		},
		{
			[4]word{0xd3917c5b, 0x55f1c407, 0x52a58a7a, 0x8f887a3b},
			[4]word{0x3e2f308c, 0xd90a8f36, 0x6ab2a923, 0x2883524c},
		},
	}

	for _, row := range table {
		qr := quarterRound3(row.in)
		if row.out != qr {
			t.Errorf("want: %x; got: %x", row.out, qr)
		}
	}
}

func BenchmarkQuarterRound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		quarterRound([4]word{0xd3917c5b, 0x55f1c407, 0x52a58a7a, 0x8f887a3b})
	}
}

func BenchmarkQuarterRound2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		quarterRound2([4]word{0xd3917c5b, 0x55f1c407, 0x52a58a7a, 0x8f887a3b})
	}
}

func BenchmarkQuarterRound3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		quarterRound3([4]word{0xd3917c5b, 0x55f1c407, 0x52a58a7a, 0x8f887a3b})
	}
}
