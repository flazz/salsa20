package salsa20

import "testing"

func TestExpansion(t *testing.T) {
	var k0, k1, n [16]byte
	for i := byte(0); i < 16; i++ {
		k0[i] = i + 1   // 1, 2, 3 …, 16
		k1[i] = i + 201 // 201, 202, 203, … 216
		n[i] = i + 101  // 101, 102, 103, … 206
	}

	// 32 bit key

	var k0k1 [32]byte
	copy(k0k1[0:16], k0[:])
	copy(k0k1[16:32], k1[:])

	got32 := expansion32(k0k1, n)
	want32 := [64]byte{
		101, 120, 112, 97, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
		13, 14, 15, 16, 110, 100, 32, 51, 101, 102, 103, 104, 105, 106, 107, 108,
		109, 110, 111, 112, 113, 114, 115, 116, 50, 45, 98, 121, 201, 202, 203, 204,
		205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 116, 101, 32, 107,
	}
	if got32 != want32 {
		t.Errorf("want: %v; got: %v", want32, got32)
	}

	gotHash32 := hash(got32)
	wantHash32 := [64]byte{
		69, 37, 68, 39, 41, 15, 107, 193, 255, 139, 122, 6, 170, 233, 217, 98,
		89, 144, 182, 106, 21, 51, 200, 65, 239, 49, 222, 34, 215, 114, 40, 126,
		104, 197, 7, 225, 197, 153, 31, 2, 102, 78, 76, 176, 84, 245, 246, 184,
		177, 160, 133, 130, 6, 72, 149, 119, 192, 195, 132, 236, 234, 103, 246, 74,
	}
	if gotHash32 != wantHash32 {
		t.Errorf("want: %v; got: %v", wantHash32, gotHash32)
	}

	// 16 bit key
	got16 := expansion16(k0, n)
	want16 := [64]byte{
		101, 120, 112, 97, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
		13, 14, 15, 16, 110, 100, 32, 49, 101, 102, 103, 104, 105, 106, 107, 108,
		109, 110, 111, 112, 113, 114, 115, 116, 54, 45, 98, 121, 1, 2, 3, 4,
		5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 116, 101, 32, 107,
	}
	if got16 != want16 {
		t.Errorf("want: %v; got: %v", want16, got16)
	}

	gotHash16 := hash(got16)
	wantHash16 := [64]byte{
		39, 173, 46, 248, 30, 200, 82, 17, 48, 67, 254, 239, 37, 18, 13, 247,
		241, 200, 61, 144, 10, 55, 50, 185, 6, 47, 246, 253, 143, 86, 187, 225,
		134, 85, 110, 246, 161, 163, 43, 235, 231, 94, 171, 51, 145, 214, 112, 29,
		14, 232, 5, 16, 151, 140, 183, 141, 171, 9, 122, 181, 104, 182, 177, 193,
	}
	if gotHash16 != wantHash16 {
		t.Errorf("want: %v; got: %v", wantHash16, gotHash16)
	}

}