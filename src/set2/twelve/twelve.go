package main

import (
	"bytes"
	"cc/util"
	"fmt"
	"strings"
)

var unKey []byte
var unString = util.Decode64("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
var altString = unString

func makeadict(in, c []byte) [][]byte {
	ret := make([][]byte, 0)
	for i := 10; i < 126; i++ {
		fill := make([]byte, 0)
		fill = append(fill, in...)
		for _, v := range c {
			fill = append(fill, v)
		}
		fill = append(fill, byte(i))
		fin := myFunc(fill, unKey)
		ret = append(ret, fin)
	}
	return ret
}

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

func gbs(key []byte) int {
	var i, score, bs int
	for i = 1; i <= 128; i++ {
		x := strings.Repeat("A", i)
		src := make([]byte, 0)
		src = append(src, []byte(x)...)
		out := myFunc(src, key)
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
	bs := gbs(unKey)
	in := []byte(strings.Repeat("A", bs-1))
	out := myFunc(in, unKey)
	dec := make([]byte, 0)
	dict := makeadict(in, dec)

	for len(dec) != len(unString) {
		for k, v := range dict {
			if bytes.Equal(v[:bs], out[:bs]) {
				c := byte(k + 10)
				dec = append(dec, c)
				fmt.Printf("%q\n", dec)
				altString = unString[len(dec):]
				makeadict(in, dec)
				out = myFunc(in, unKey)
				break
			}
		}
	}
	fmt.Println(string(dec))
}
