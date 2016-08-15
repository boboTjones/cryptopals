/*
3. Single-character XOR Cipher

The hex encoded string:

      1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736

... has been XOR'd against a single character. Find the key, decrypt
the message.

Write code to do this for you.

Here's one way:

a. Find a large sample of English text. Something from Project
Gutenberg should do nicely. Use it to generate a character frequency
map.

b. Evaluate each potential key by scoring the resulting plaintext
against the frquency map. The key with the best score is your match.

*/
package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
    "sort"
)

type Result struct {
    key byte
	text  []byte
	total int
}

type Results []*Result

func (r Results) Len() int { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type sortByTotal struct { Results }

func (s sortByTotal) Less(i, j int) bool { return s.Results[i].total < s.Results[j].total }

func someText() []byte {
	r, err := http.Get("http://www.gutenberg.org/files/205/205-h/205-h.htm")
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	return b
}

func mapMaker(text []byte) map[byte]int {
	x := make(map[byte]int)
	for _, t := range text {
		x[t] += 1
	}
	return x
}

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

	results := Results{}
	cfm := mapMaker(someText())
	w := deFunk("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")

	for z := 32; z < 126; z++ {
		results = append(results, &Result{uint8(z), deXor(w, uint8(z)), 0})
	}

	for k, v := range results {
		q := mapMaker(v.text)
		y := 0
		for l, _ := range q {
			y += cfm[l]
			results[k].total = y
		}
	}
    
    sort.Sort(sortByTotal{results})
    
    for _, v := range(results) {
        fmt.Printf("%s (%v): %d\n", string(v.key), v.key, v.total)
    }
}
