package salsa20

type word uint32

func quarterRound(y [4]word) [4]word {
	var z [4]word
	var u word

	// z1 = y1 ⊕ ((y0 + y3) <<< 7),
	u = y[0] + y[3]
	z[1] = y[1] ^ (u<<7 | u>>(32-7))

	// z2 = y2 ⊕ ((z1 + y0) <<< 9),
	u = z[1] + y[0]
	z[2] = y[2] ^ (u<<9 | u>>(32-9))

	// z3 = y3 ⊕ ((z2 + z1) <<< 13),
	u = z[2] + z[1]
	z[3] = y[3] ^ (u<<13 | u>>(32-13))

	// z0 = y0 ⊕ ((z3 + z2) <<< 18).
	u = z[3] + z[2]
	z[0] = y[0] ^ (u<<18 | u>>(32-18))

	return z
}

func rowRound(y [16]word) [16]word {
	var z [16]word
	var q [4]word

	// (z0, z1, z2, z3) = quarterround(y0, y1, y2, y3),
	q = quarterRound([4]word{y[0], y[1], y[2], y[3]})
	z[0], z[1], z[2], z[3] = q[0], q[1], q[2], q[3]

	// (z5, z6, z7, z4) = quarterround(y5, y6, y7, y4),
	q = quarterRound([4]word{y[5], y[6], y[7], y[4]})
	z[5], z[6], z[7], z[4] = q[0], q[1], q[2], q[3]

	// (z10, z11, z8, z9) = quarterround(y10, y11, y8, y9),
	q = quarterRound([4]word{y[10], y[11], y[8], y[9]})
	z[10], z[11], z[8], z[9] = q[0], q[1], q[2], q[3]

	// (z15, z12, z13, z14) = quarterround(y15, y12, y13, y14).
	q = quarterRound([4]word{y[15], y[12], y[13], y[14]})
	z[15], z[12], z[13], z[14] = q[0], q[1], q[2], q[3]

	return z
}

func columnRound(x [16]word) [16]word {
	var y [16]word
	/*
	   (y0, y4, y8, y12, y1, y5, y9, y13, y2, y6, y10, y14, y3, y7, y11, y15) =
	   rowround(x0, x4, x8, x12, x1, x5, x9, x13, x2, x6, x10, x14, x3, x7, x11, x15).
	*/
	r := rowRound([16]word{
		x[0], x[4], x[8], x[12],
		x[1], x[5], x[9], x[13],
		x[2], x[6], x[10], x[14],
		x[3], x[7], x[11], x[15],
	})

	y[0], y[4], y[8], y[12],
		y[1], y[5], y[9], y[13],
		y[2], y[6], y[10], y[14],
		y[3], y[7], y[11], y[15] =
		r[0], r[1], r[2], r[3],
		r[4], r[5], r[6], r[7],
		r[8], r[9], r[10], r[11],
		r[12], r[13], r[14], r[15]

	return y
}

func doubleRound(x [16]word) [16]word {
	return rowRound(columnRound(x))
}

func littleEndian(b [4]byte) word {
	// littleendian(b) = b0 + 2^8*b1 + 2^16*b2 + 2^24*b3.
	w := word(b[0])
	w += word(b[1]) << 8
	w += word(b[2]) << 16
	w += word(b[3]) << 24

	return w
}

func littleEndianInv(w word) [4]byte {
	var b [4]byte

	b[0] = byte(w)
	b[1] = byte(w >> 8)
	b[2] = byte(w >> 16)
	b[3] = byte(w >> 24)

	return b
}

func hash(b [64]byte) [64]byte {
	var x [16]word
	for i := 0; i < 16; i++ {
		s := 4 * i
		x[i] = littleEndian([4]byte{b[s+0], b[s+1], b[s+2], b[s+3]})
	}

	// (z0, z1, . . . , z15) = doubleround10(x0, x1, . . . , x15).
	z := x
	for i := 0; i < 10; i++ {
		z = doubleRound(z)
	}

	var r [64]byte
	for i := 0; i < 16; i++ {
		s := 4 * i
		b := littleEndianInv(z[i] + x[i])
		for j := 0; j < 4; j++ {
			r[s+j] = b[j]
		}
	}

	return r
}

var (
	sigma = [4][4]byte{
		[4]byte{101, 120, 112, 97},
		[4]byte{110, 100, 32, 51},
		[4]byte{50, 45, 98, 121},
		[4]byte{116, 101, 32, 107},
	}

	tau = [4][4]byte{
		[4]byte{101, 120, 112, 97},
		[4]byte{110, 100, 32, 49},
		[4]byte{54, 45, 98, 121},
		[4]byte{116, 101, 32, 107},
	}
)

// Salsa20(σ0, k0, σ1, n, σ2, k1, σ3)
func expansion32(k [32]byte, n [16]byte) [64]byte {
	var s [64]byte

	copy(s[0:4], sigma[0][:])
	copy(s[4:20], k[0:16])
	copy(s[20:24], sigma[1][:])
	copy(s[24:40], n[:])
	copy(s[40:44], sigma[2][:])
	copy(s[44:60], k[16:32])
	copy(s[60:64], sigma[3][:])

	return s
}

// Salsa20(τ0, k, τ1, n, τ2, k, τ3)
func expansion16(k [16]byte, n [16]byte) [64]byte {
	var s [64]byte

	copy(s[0:4], tau[0][:])
	copy(s[4:20], k[:])
	copy(s[20:24], tau[1][:])
	copy(s[24:40], n[:])
	copy(s[40:44], tau[2][:])
	copy(s[44:60], k[:])
	copy(s[60:64], tau[3][:])

	return s
}
