// Byte-at-a-time ECB decryption (Harder)
package main

import (
	"bytes"
	"cc/util"
	"fmt"
	"math/rand"
	"time"
)

func mkSrc(rp, ac, tb []byte) []byte {
	src := make([]byte, 0)
	src = append(src, rp...)
	src = append(src, ac...)
	src = append(src, tb...)
	return src
}

func main() {
	key := util.RandString(16)
	//Now generate a random count of random bytes and prepend this string to every plaintext.
	// AES-128-ECB(random-prefix || attacker-controlled || target-bytes, random-key)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// random-prefix
	rp := util.RandString(r1.Intn(33) + 2)
	// attacker-controlled
	ac := bytes.Repeat([]byte("A"), 15)
	// target-bytes
	tb := util.Decode64("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
	spare := tb

	in := mkSrc(rp, ac, spare)
	ecb := util.NewECB(key)
	out := ecb.Orakkl(in)

	dec := make([]byte, 0)
	dict := ecb.Dict(rp, ac, tb, dec)
	fmt.Println("WTF ", dict[len(dict)-1:])

	for len(dec) != len(tb) {
		for k, v := range dict {
			fmt.Printf("dict\t%v\n", v[:16])
			fmt.Printf("out \t%v\n", out[:16])
			if bytes.Equal(v[:16], out[:16]) {
				dec = append(dec, byte(k+10))
				fmt.Printf("Found %q\n", dec)
				spare = spare[len(dec):]
				in = mkSrc(rp, ac, spare)
				dict = ecb.Dict(rp, ac, tb, dec)
				out = ecb.Orakkl(in)
				break
			}
		}
	}
	fmt.Println(string(dec))
}
