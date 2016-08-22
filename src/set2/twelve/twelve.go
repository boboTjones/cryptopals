package main

import (
	"cc/util"
	"fmt"
	"strings"
)

var unKey []byte

func makeadict(key []byte, size int) [][]byte {
	ret := make([][]byte, 0)
	for i := 66; i < 113; i++ {
		fill := []byte(strings.Repeat("A", size))
		fill = append(fill, byte(i))
		fmt.Println(fill)
		fin := myFunc(fill, key)
		ret = append(ret, fin)
	}
	return ret
}

func detect(src []byte) int {
	return util.Compare(src, 16)
}

func myFunc(src, key []byte) []byte {
	// unknown-string
	us := []byte("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
	src = append(src, us...)

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

func lite(src, key []byte) []byte {
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
	fmt.Printf("Found block size %d\n", bs)
	in := []byte(strings.Repeat("A", bs-1))
	out := myFunc(in, unKey)

	for i := bs - 1; i < 184; i++ {
		dict := makeadict(unKey, i)
		for k, v := range dict {
			e := i + 1
			fmt.Printf("%q\n", (k + 66))
			fmt.Println("DI  ", v[:e])
			fmt.Println("OU  ", out[:e])
			//if bytes.Equal(v[:e], out[:e]) {
			//	fmt.Println("MOOSE")
			//}
		}
	}

}
