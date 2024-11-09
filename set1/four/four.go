/*
4. Detect single-character XOR

One of the 60-character strings at:

	https://gist.github.com/3132713

has been encrypted by single-character XOR. Find it. (Your code from
#3 should help.)
*/
package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/bobotjones/cryptopals/util"
)

var fileName string

type Result struct {
	original  []byte
	lineNo    int
	key       byte
	decrypted []byte
	total     int
}

type Results []*Result

func (r Results) Len() int      { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type sortByTotal struct{ Results }

func (s sortByTotal) Less(i, j int) bool { return s.Results[i].total < s.Results[j].total }

func deFunk(str string) []byte {
	foo := make([]byte, len(str)/2)
	_, err := hex.Decode(foo, []byte(str))
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	return foo
}

func deXor(str []byte, char byte) []byte {
	var x []byte
	for i := 0; i < len(str); i++ {
		x = append(x, str[i]^char)
	}
	return x
}

func init() {
	flag.StringVar(&fileName, "f", fileName, "File to decrypt")
}

func main() {
	flag.Parse()
	data := bytes.NewBuffer(make([]byte, 0))

	if fileName != "" {
		d, err := util.SlurpFromFile(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data.Write(d)
	} else {
		fmt.Println("need input")
		os.Exit(1)
	}

	nl := uint8(10)
	text := util.FetchStuff("http://www.gutenberg.org/files/205/205-h/205-h.htm")
	leMap := util.MapMaker(text)
	cMap := util.CharCount(text)
	results := Results{}

	lines := bytes.Split(data.Bytes(), []byte{nl})

	for i, l := range lines {
		for _, z := range leMap {
			key := uint8(uint8(z.Char))
			orig := deFunk(string(l))
			decr := deXor(orig, key)
			total := 0
			for _, c := range decr {
				rr := cMap[c]
				total += rr
			}
			// original text, line number, key, decrypted text, total, score
			results = append(results, &Result{orig, i, key, decr, total})
		}
	}

	sort.Sort(sortByTotal{results})

	for _, result := range results {
		fmt.Printf("Score for line %v with key %v(%s): %v\n%s\n", result.lineNo, result.key, string(result.key), result.total, hex.Dump(result.decrypted))
	}
}
