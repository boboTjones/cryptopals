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

func makeadict(in []byte, c []byte) [][]byte {
	//fmt.Printf("Got %d bytes\n", len(in))
	ret := make([][]byte, 0)
	for i := 33; i < 126; i++ {
		fill := make([]byte, 0)
		fill = append(fill, in...)
		for _, v := range c {
			fill = append(fill, v)
		}
		fill = append(fill, byte(i))
		//fmt.Printf("%d in: %q\n", len(fill), fill)
		fin := myFunc(fill, unKey)
		ret = append(ret, fin)
	}
	return ret
}

func myFunc(src, key []byte) []byte {
	//fmt.Println("IN  ", string(src))
	// unknown-string

	src = append(src, altString...)
	// pad
	fmt.Println(string(src))
	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	//fmt.Printf("%q\n", src)
	out := make([]byte, len(src))
	ecb := util.NewECB(key)
	ecb.Encrypt(out, src)
	//fmt.Printf("inter\t%v\n", out[:16])
	return out
}

func gbs(key []byte) int {
	var i, score, bs int
	for i = 1; i <= 128; i++ {
		x := strings.Repeat("A", i)
		//fmt.Println("Input %d\t%s\n", i, x)
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
	//fmt.Printf("Found block size %d\n", bs)
	o := bs - 1

	in := []byte(strings.Repeat("A", o))
	out := myFunc(in, unKey)
	dec := make([]byte, 0)
	dict := makeadict(in, dec)

	for len(dec) != len(unString) {
		for k, v := range dict {
			//fmt.Printf("%q\n", (k + 33))
			//fmt.Printf("%d\t%v\n", len(v[:bs]), v[:bs])
			//fmt.Printf("%d\t%v\n", len(out[:bs]), out[:bs])
			if bytes.Equal(v[:bs], out[:bs]) {
				c := byte(k + 33)
				dec = append(dec, c)
				fmt.Printf("MOOSE\t%q\n", dec)
				//o--
				altString = unString[len(dec):]
				in = []byte(strings.Repeat("A", o))
				makeadict(in, dec)
				out = myFunc(in, unKey)
				break
			}
		}
	}
	fmt.Println(string(dec))
}
