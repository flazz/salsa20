package salsa20

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/flazz/salsa20/xor"
)

func TestEncrypt(t *testing.T) {
	k32 := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		201, 202, 203, 204, 205, 206, 207, 208,
		209, 210, 211, 212, 213, 214, 215, 216,
	}
	nonce := [8]byte{
		101, 102, 103, 104, 105, 106, 107, 108,
	}

	cip16, err := NewCipher(k32[:16], nonce)
	if err != nil {
		fmt.Print(err)
		return
	}

	cip32, err := NewCipher(k32, nonce)
	if err != nil {
		fmt.Print(err)
		return
	}

	ciphers := []Cipher{cip16, cip32}

	messages := []string{
		"foo bar baz",
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
		eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
		minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
		ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
		voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur
		sint occaecat cupidatat non proident, sunt in culpa qui officia
		deserunt mollit anim id est laborum.`,
	}

	ciphers = ciphers[0:1]
	messages = messages[0:1]

	for _, cip := range ciphers {
		for _, msg := range messages {
			plainText := []byte(msg)

			encr := xor.NewReader(cip.NewReader(), bytes.NewReader(plainText))

			cipherText, err := ioutil.ReadAll(encr)
			if err != nil {
				t.Error(err)
			}
			t.Log(cipherText)
		}
	}
}
