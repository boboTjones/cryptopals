/*
7. AES in ECB Mode

The Base64-encoded content at the following location:

   https://gist.github.com/3132853

Has been encrypted via AES-128 in ECB mode under the key

   "YELLOW SUBMARINE".

(I like "YELLOW SUBMARINE" because it's exactly 16 bytes long).

Decrypt it.

Easiest way:

Use OpenSSL::Cipher and give it AES-128-ECB as the cipher.
*/

package main

import (
	"bytes"
	"crypto/aes"
	"flag"
	"fmt"
	"github.com/bobotjones/cryptopals/util"
	"os"
)

var fileName string

func init() {
	flag.StringVar(&fileName, "f", fileName, "File to decrypt")
}

func main() {
	flag.Parse()

	data := bytes.NewBuffer(make([]byte, 0))

	if fileName != "" {
		d := util.SlurpFromFile(fileName)
		encrypted := util.Decode64(string(d))
		data.Write(encrypted)
	} else {
		fmt.Println("need input")
		os.Exit(1)
	}

	encKey := []byte("YELLOW SUBMARINE")

	c, _ := aes.NewCipher(encKey)

	src := data.Bytes()
	out := make([]byte, len(src))
	dst := out
	bs := 16

	if len(src)%bs != 0 {
		fmt.Println("not blocks")
		os.Exit(1)
	}

	for len(src) > 0 {
		c.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}

	fmt.Printf("Output:\n%s", string(out))
}
