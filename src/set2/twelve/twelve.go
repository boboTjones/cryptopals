package main

import (
	"bytes"
	"cc/util"
	"fmt"
	"strings"
)

var unKey []byte

func makeadict(size int, c ...byte) [][]byte {
	ret := make([][]byte, 0)
	for i := 33; i < 126; i++ {
		fill := []byte(strings.Repeat("A", size))
		if c != nil {
			fill = append(fill, c[0])
		}
		fmt.Println("DI a  ", fill)
		fill = append(fill, byte(i))
		fin := lite(fill, unKey)
		fmt.Println("DI b  ", fin[:size])
		ret = append(ret, fin)
	}
	return ret
}

func detect(src []byte) int {
	return util.Compare(src, 16)
}

func myFunc(src, key []byte) []byte {
	//fmt.Println("IN  ", string(src))
	// unknown-string
	us := util.Decode64("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
	src = append(src, us...)
	// pad
	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	fmt.Printf("%q\n", src)
	out := make([]byte, len(src))
	ecb := util.NewECB(key)
	ecb.Encrypt(out, src)
	//fmt.Printf("inter\t%v\n", out[:16])
	return out
}

func lite(src, key []byte) []byte {
	// pad
	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	fmt.Printf("%q\n", src)
	out := make([]byte, len(src))
	ecb := util.NewECB(key)
	ecb.Encrypt(out, src)
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
	fmt.Println(string(unKey))
	fmt.Println(lite([]byte("AAAAAAAAAAAAAAARo"), unKey))

	bs := gbs(unKey)
	fmt.Printf("Found block size %d\n", bs)

	in := []byte(strings.Repeat("A", bs-1))
	out := myFunc(in, unKey)
	//fmt.Println(util.Chunk(out, 16))

	dict := makeadict(bs - 1)
	e := bs
	for {
		for k, v := range dict {
			fmt.Printf("%q\n", (k + 33))

			fmt.Println("OU  ", out[:e])
			if bytes.Equal(v[:e], out) {
				c := byte(k + 33)
				fmt.Printf("MOOSE\t%q\n", c)
				e++
				makeadict(bs-1, c)
				break
			}
		}
	}
}
