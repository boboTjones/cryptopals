package util

import (
	"crypto/aes"
	"crypto/cipher"
	//"fmt"
)

type CBC struct {
	blockSize int
	iv        []byte
	thing     cipher.Block
}

func NewCBC(iv, key []byte) *CBC {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &CBC{
		blockSize: 16,
		iv:        iv,
		thing:     c,
	}
}

func (c *CBC) Encrypt(dst, src []byte) {
	iv := c.iv
	s := 0
	e := c.blockSize

	for i := 0; i < len(src)/c.blockSize; i++ {
		//fmt.Printf("encrypting block %d-%d\n", s, e)
		Xor(dst[s:e], src[s:e], iv)
		c.thing.Encrypt(dst[s:e], dst[s:e])
		iv = dst[s:e]
		s += c.blockSize
		e = s + c.blockSize
	}
}

func (c *CBC) Decrypt(dst, src []byte) {
	var iv []byte
	e := len(src)
	for s := e - c.blockSize; s >= 0; s -= c.blockSize {
		//fmt.Printf("decrypting block %d-%d\n", s, e)
		out := dst[s:e]
		in := src[s:e]
		c.thing.Decrypt(out, in)

		if s == 0 {
			iv = c.iv
		} else {
			iv = src[(s - c.blockSize):s]
		}
		Xor(out, out, iv)
		e = s
	}
}
