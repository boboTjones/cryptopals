// Byte-at-a-time ECB decryption (Harder)
package main

import (
	"bytes"
	"cc/util"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var unKey []byte
var unString = util.Decode64("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
var altString = unString

func makeadict(in, c []byte) [][]byte {
	ret := make([][]byte, 0)
	for i := 10; i < 126; i++ {
		fill := make([]byte, 0)
		fill = append(fill, in...)
		fill = append(fill, c...)
		fill = append(fill, byte(i))
		fin := myFunc(fill)
		ret = append(ret, fin)
	}
	return ret
}

func myFunc(src []byte) []byte {
	// random-prefix || attacker-controlled || target-bytes
	src = append(src, altString...)
	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	out := make([]byte, len(src))
	ecb := util.NewECB(unKey)
	ecb.Encrypt(out, src)
	return out
}

func main() {
	//Now generate a random count of random bytes and prepend this string to every plaintext.
	// AES-128-ECB(random-prefix || attacker-controlled || target-bytes, random-key)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// random-prefix
	rp := util.RandString(r1.Intn(33) + 2)

	// attacker-controlled, padding out random to full block
	x := 16 - math.Mod(float64(len(rp)), 16)
	ac := bytes.Repeat([]byte("A"), int(x)-1)

	// What if I change the block size?
	bs := len(ac) + len(rp) + 1
	fmt.Printf("Using block size %d\n", bs)
	unKey = util.RandString(bs)

	// random-prefix || attacker-controlled
	in := append(rp, ac...)
	out := myFunc(in)

	dec := make([]byte, 0)
	dict := makeadict(in, dec)

	for len(dec) != len(unString) {
		for k, v := range dict {
			if bytes.Equal(v[:bs], out[:bs]) {
				dec = append(dec, byte(k+10))
				fmt.Printf("Found %q\n", dec)
				altString = unString[len(dec):]
				out = myFunc(in)
				break
			}
		}
	}
	fmt.Println(string(dec))
}
