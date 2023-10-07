/*
CBC mode is a block cipher mode that allows us to encrypt
irregularly-sized messages, despite the fact that a block cipher
natively only transforms individual blocks.

In CBC mode, each ciphertext block is added to the next plaintext
block before the next call to the cipher core.

The first plaintext block, which has no associated previous ciphertext
block, is added to a "fake 0th ciphertext block" called the
initialization vector, or IV.

Implement CBC mode by hand by taking the ECB function you wrote
earlier, making it encrypt instead of decrypt (verify this by
decrypting whatever you encrypt to test), and using your XOR function
from the previous exercise to combine them.

The file here is intelligible (somewhat) when CBC decrypted against
"YELLOW SUBMARINE" with an IV of all ASCII 0 (\x00\x00\x00 &c)
*/
package main

import (
	"bytes"
	"fmt"
	"util"
)

func getEncData() []byte {
	data := &bytes.Buffer{}
	data.Write(util.Decode64(string(util.SlurpFromFile("/Users/erin/codebase/cryptochallenges/src/set2/ten/10.txt"))))
	return data.Bytes()
}

func getPlainText(filePath string) []byte {
	data := &bytes.Buffer{}
	data.Write(util.SlurpFromFile(filePath))
	return data.Bytes()
}

func main() {
	iv := []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	k := []byte("YELLOW SUBMARINE")
	cbc := util.NewCBC(iv, k)
	src := getPlainText("/Users/erin/encme.txt")

	if len(src)%16 != 0 {
		pddr := util.NewPadder(16)
		pddr.Data.Write(src)
		pddr.Padfoot()
		src = pddr.Data.Bytes()
	}

	fmt.Printf("original (%d bytes): %#v\n", len(src), src)

	dst := make([]byte, len(src))
	cbc.Encrypt(dst, src)
	fmt.Printf("encrypted (%d bytes): %#v\n", len(dst), dst)

	dst2 := make([]byte, len(dst))
	cbc.Decrypt(dst2, dst)
	fmt.Println(string(dst2))

	//	src3 := getEncData()
	//	cbc3 := util.NewCBC(iv, k)
	//	dst3 := make([]byte, len(src3))
	//	cbc3.Decrypt(dst3, src3)
	//	fmt.Println(string(dst3))
}
