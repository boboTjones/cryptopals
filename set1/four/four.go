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
	"cfm"
	"encoding/hex"
	"fmt"
	"sort"
)

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

func main() {
	nl := uint8(10)
	//encryptedTextURL := "https://gist.github.com/tqbf/3132713/raw/40da378d42026a0731ee1cd0b2bd50f66aabac5b/gistfile1.txt"
	// encryptedText := cfm.FetchStuff(encryptedTextURL)
	encryptedText := cfm.Slurp("/Users/erin/codebase/cryptochallenges/src/cc/four/gistfile1.txt")
	text := cfm.FetchStuff("http://www.gutenberg.org/files/205/205-h/205-h.htm")
	leMap := cfm.MapMaker(text)
	cMap := cfm.CharCount(text)
	results := Results{}

	lines := bytes.Split(encryptedText, []byte{nl})

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
		fmt.Printf("Score for line %v with key %v(%s): %v\nText: %s\n\n", result.lineNo, result.key, string(result.key), result.total, string(result.decrypted))
	}
}
