/*
Stuff I keep using.
*/

package util

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

func HexStringToByte(str string) []byte {
	foo := make([]byte, len(str)/2)
	_, err := hex.Decode(foo, []byte(str))
	if err != nil {
		fmt.Println("Something bad happened.")
		os.Exit(1)
	}
	return foo
}

func DeXor(str []byte, char byte) []byte {
	var x []byte
	for i := 0; i < len(str); i++ {
		x = append(x, str[i]^char)
	}
	return x
}

func Xor(dst, in, iv []byte) int {
	n := len(in)
	if len(iv) < n {
		n = len(iv)
	}
	for i := 0; i < n; i++ {
		dst[i] = in[i] ^ iv[i]
	}
	return n
}

func SlurpFromFile(filePath string) []byte {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
		os.Exit(1)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
		os.Exit(1)
	}
	return data
}

func SlurpFromURL(t string) []byte {
	r, err := http.Get(t)
	if err != nil {
		fmt.Println("Something bad happened: %v", err)
		os.Exit(1)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
		os.Exit(1)
	}
	return b
}

func Decode64(str string) []byte {
	t, _ := base64.StdEncoding.DecodeString(str)
	return t
}

func Chunk(blob []byte, csize int) [][]byte {
	var fin = make([][]byte, 0)
	r := int((math.Mod(float64(len(blob)), float64(csize))))
	its := (len(blob) - r) / csize
	s := 0
	e := 0
	for i := 0; i < its; i++ {
		s = i * csize
		e = s + csize
		fin = append(fin, blob[s:e])
	}
	if r != 0 {
		fin = append(fin, blob[s+csize:])
	}
	return fin
}

// Courtesy of Andy Schmitz.

func AndyChunk(blob []byte, csize int) [][]byte {
	var fin = make([][]byte, 0)
	for i := 0; i < len(blob); i += csize {
		e := i + csize
		if e > len(blob) {
			e = len(blob)
		}
		fin = append(fin, blob[i:e])
	}
	return fin
}
