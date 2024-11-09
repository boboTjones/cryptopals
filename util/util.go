/*
Stuff I keep using.
*/

package util

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
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

func SlurpFromFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
		return nil, err
	}
	return data, nil
}

func SlurpFromURL(t string) ([]byte, error) {
	r, err := http.Get(t)
	if err != nil {
		fmt.Println("Something bad happened: %v", err)
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
		return nil, err
	}
	return b, nil
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

// Courtesy of aschmitz.

func AshChunk(blob []byte, csize int) [][]byte {
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

func RandString(n int) []byte {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789987654321123456789abcdefghijklmnopqrstuvwxyz")
	ret := make([]byte, n)
	x := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, x)
	if err != nil {
		fmt.Println("rand: ", err)
	}
	max := byte(len(chars))
	for i, c := range x {
		ret[i] = chars[c%max]
	}
	return ret
}

func Compare(c []byte, n int) int {
	var ret int
	chunks := Chunk(c, n)
	x := len(chunks) - 1
	for i := 0; i < x; i++ {
		if bytes.Equal(chunks[i], chunks[i+1]) {
			//fmt.Printf("MOOSE: %d and %d match\n", i, i+1)
			ret++
		} else {
			//fmt.Printf("%d and %d do not match\n", i, i+1)
		}
		//fmt.Printf("%v\n", chunks[i])
		//fmt.Printf("%v\n\n", chunks[i+1])
	}
	return ret
}

func ReComp(c []byte, n int) int {
	var ret int
	fmt.Println(c)
	src := Chunk(c, n)

	for x := 0; x < len(src)-1; x++ {
		hi := src[x]
		//fmt.Printf("Comparing %d\t%v...\n", x, hi)
		for y, s := range src {
			if x == y {
				continue
			}
			//fmt.Printf("...to %d\t\t%v\n", y, s)
			for i, v := range s {
				if hi[i] == v {
					ret++
				}
			}
		}
	}
	fmt.Println(ret)
	return ret
}
