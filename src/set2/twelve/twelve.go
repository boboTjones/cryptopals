package main

import (
	"cc/util"
	"fmt"
)

var unKey []byte

func main() {
	unKey = util.RandString(16)
	mysteryText := "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
	src := make([]byte, 160)
	for i, _ := range src {
		src[i] = byte(0x41)
	}
	src = append(src, util.Decode64(mysteryText)...)
	pddr := util.NewPadder(16)
	pddr.Data.Write(src)
	pddr.Padfoot()
	in := pddr.Data.Bytes()
	out := make([]byte, len(in))
	ecb := util.NewECB(unKey)
	ecb.Encrypt(out, in)
	fmt.Println(out)

}
