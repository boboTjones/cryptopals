package main

import (
	"bytes"
	"fmt"
	"strings"
	"util"
)

var unKey []byte
var unString = util.Decode64("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
var altString = unString

func makeadict(in, c []byte) [][]byte {
	ret := make([][]byte, 0)
	for i := 10; i < 126; i++ {
		fmt.Printf("CHAR\t%q\n", i)

		fill := make([]byte, 0)
		fill = append(fill, in...)
		for _, v := range c {
			fill = append(fill, v)
		}
		fill = append(fill, byte(i))
		fin := myFunc(fill)
		ret = append(ret, fin)
	}
	return ret
}

func myFunc(src []byte) []byte {
	// random-prefix || attacker-controlled || target-bytes
	src = append(src, altString...)
	// pad
	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	fmt.Printf("PADD\t%q\n", src)
	out := make([]byte, len(src))
	ecb := util.NewECB(unKey)
	ecb.Encrypt(out, src)
	fmt.Printf("ENCR\t%v\n", out)
	return out
}

func gbs() int {
	var i, score, bs int
	for i = 1; i <= 128; i++ {
		x := strings.Repeat("A", i)
		src := make([]byte, 0)
		src = append(src, []byte(x)...)
		out := myFunc(src)
		score = util.Compare(out, 16)
		if score > 0 {
			bs = i / 2
			break
		}
	}
	return bs
}

func main() {
	unKey = util.RandString(16)
	bs := gbs()
	ayes := []byte(strings.Repeat("A", bs-1))
	out := myFunc(ayes)
	dec := make([]byte, 0)
	dict := makeadict(ayes, dec)

	for len(dec) != len(unString) {
		for k, v := range dict {
			if bytes.Equal(v[:bs], out[:bs]) {
				c := byte(k + 10)
				dec = append(dec, c)
				fmt.Printf("Found %q\n", dec)
				altString = unString[len(dec):]
				out = myFunc(ayes)
				break
			}
		}
	}
	fmt.Println(string(dec))
}
