package util

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type ECB struct {
	thing     cipher.Block
	blockSize int
	key       []byte
}

func NewECB(key []byte) *ECB {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &ECB{thing: c, blockSize: 16}
}

func (e *ECB) Encrypt(dst, src []byte) {
	for len(src) > 0 {
		e.thing.Encrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

func (e *ECB) Decrypt(dst, src []byte) {
	for len(src) > 0 {
		e.thing.Decrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

func (e *ECB) Blocksize() int {
	return e.blockSize
}

func (e *ECB) Dict(rp, ac, tb, c []byte) [][]byte {
	fmt.Println("Make dict")
	ret := make([][]byte, 0)
	for i := 10; i < 126; i++ {
		fill := make([]byte, 0)
		fill = append(fill, rp...)
		fill = append(fill, ac...)
		fill = append(fill, c...)
		fill = append(fill, byte(i))
		fill = append(fill, tb...)
		fmt.Printf("Fill\t%q\n", fill)
		fin := e.Orakkl(fill)
		ret = append(ret, fin)
	}
	return ret
}

func (e *ECB) Orakkl(src []byte) []byte {
	fmt.Println("Oracling")
	if len(src)%16 != 0 {
		pddr := NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}
	fmt.Printf("Padd\t%q\n", src)
	out := make([]byte, len(src))
	e.thing.Encrypt(out, src)
	return out
}
