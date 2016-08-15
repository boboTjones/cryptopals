/*
2. Fixed XOR

Write a function that takes two equal-length buffers and produces
their XOR sum.

The string:

 1c0111001f010100061a024b53535009181c

... after hex decoding, when xor'd against:

 686974207468652062756c6c277320657965

... should produce:

 746865206b696420646f6e277420706c6179

*/

package main

import (
	"encoding/hex"
	"fmt"
)

func deFunk(str string) []byte {
	foo := make([]byte, len(str)/2)
	_, err := hex.Decode(foo, []byte(str))
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	return foo
}

func main() {

	strone := deFunk("1c0111001f010100061a024b53535009181c")
	strtwo := deFunk("686974207468652062756c6c277320657965")
	var buf []byte
    
	for i := 0; i < len(strone); i++ {
		buf = append(buf, (strone[i] ^ strtwo[i]))
	}
    
	fmt.Printf("%x\n", buf)
}
