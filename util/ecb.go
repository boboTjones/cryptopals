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
