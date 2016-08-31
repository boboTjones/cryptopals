// Byte-at-a-time ECB decryption (Harder)
package main

import (
	"fmt"
)

// Take your oracle function from #12.

func myFunc(src, key []byte) []byte {
	// unknown-string
	src = append(src, altString...)
	// pad
	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	out := make([]byte, len(src))
	ecb := util.NewECB(key)
	ecb.Encrypt(out, src)
	return out
}

func main() {
	//Now generate a random count of random bytes and prepend this string to every plaintext.
	// AES-128-ECB(random-prefix || attacker-controlled || target-bytes, random-key)
	fmt.Println("helu")
}
