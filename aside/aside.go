package main

import (
	"bytes"
	"fmt"
	"github.com/bobotjones/cryptopals/util"
	"strings"
)

var unString = util.Decode64("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
var altString = unString

func makeadict(in, c, key, iv []byte) [][]byte {
	ret := make([][]byte, 0)
	for i := 0; i < 127; i++ {
		fill := make([]byte, 0)
		fill = append(fill, in...)
		for _, v := range c {
			fill = append(fill, v)
		}
		fill = append(fill, byte(i))
		fin := myFunc(fill, key, iv)
		ret = append(ret, fin)
	}
	return ret
}

func myFunc(src, key, iv []byte) []byte {
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
	cbc := util.NewCBC(key, iv)
	cbc.BadEncrypt(out, src)
	return out
}

func gbs(key, iv []byte) int {
	var i, score, bs int
	for i = 1; i <= 128; i++ {
		x := strings.Repeat("A", i)
		src := make([]byte, 0)
		src = append(src, []byte(x)...)
		out := myFunc(src, key, iv)
		score = util.Compare(out, 16)
		if score > 0 {
			bs = i / 2
			break
		}
	}
	return bs
}

func main() {
	key := util.RandString(16)
	iv := util.RandString(16)

	bs := gbs(key, iv)
	fmt.Printf("block size is %d\n", bs)

	src := []byte(strings.Repeat("A", 15))
	dst := myFunc(src, key, iv)
	dec := make([]byte, 0)
	dict := makeadict(src, dec, key, iv)

	for len(dec) != len(unString) {
		for k, v := range dict {
			if bytes.Equal(v[:bs], dst[:bs]) {
				c := byte(k)
				dec = append(dec, c)
				fmt.Printf("%q\n", dec)
				altString = unString[len(dec):]
				makeadict(src, dec, key, iv)
				dst = myFunc(src, key, iv)
				break
			}
		}
	}
	fmt.Println(string(dec))
}
