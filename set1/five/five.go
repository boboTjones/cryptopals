/*
5. Repeating-key XOR Cipher

Write the code to encrypt the string:

  Burning 'em, if you ain't quick and nimble
  I go crazy when I hear a cymbal

Under the key "ICE", using repeating-key XOR. It should come out to:

0b 36 37272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f
0b 01 07272d2b2e24222c26206963656963653c3630202a2c3d37313c363022282e272d2b272d2b2b212743494f2e24222a202633393f3e3432272d2b696365282224282224303a3c282224

0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27

Encrypt a bunch of stuff using your repeating-key XOR function. Get a feel for it.
*/

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/bobotjones/cryptopals/util"
)

var plainText, fileName, keyString string

func init() {
	flag.StringVar(&plainText, "t", plainText, "String to encrypt")
	flag.StringVar(&fileName, "f", fileName, "File to encrypt")
	flag.StringVar(&keyString, "k", keyString, "Key value")
}

func main() {
	var data []byte
	var result []byte
	var key []byte

	flag.Parse()
	switch {
	case fileName != "":
		data = util.SlurpFromFile(fileName)
	case plainText != "":
		data = []byte(plainText)
	default:
		data = []byte("Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal")
	}

	if keyString == "" {
		key = []byte("ICE")
	} else {
		key = []byte(keyString)
	}

	for i := 0; i < len(data); i++ {
		result = append(result, (data[i] ^ key[i%len(key)]))
	}

	fmt.Println(hex.EncodeToString(result))
}
