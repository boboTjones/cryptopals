package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type Result struct {
	key   byte
	text  []byte
	score float64 // Changed from int to float64 for relative frequencies
	total int
}

type Results []*Result

func (r Results) Len() int      { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type sortByScore struct{ Results }

func (s sortByScore) Less(i, j int) bool { return s.Results[i].score < s.Results[j].score }

// Modified to return both the frequency map and total count
func mapMaker(text []byte) (map[byte]float64, int) {
	freqs := make(map[byte]float64)
	total := 0
	for _, t := range text {
		freqs[t]++
		total++
	}

	// Convert to relative frequencies
	for k := range freqs {
		freqs[k] = freqs[k] / float64(total)
	}
	return freqs, total
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

func getSomeText() []byte {
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

func main() {
	results := Results{}

	// Get reference text and calculate relative frequencies
	rFreqs, _ := mapMaker(getSomeText())

	// The target ciphertext
	w := deFunk("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")

	// Try each possible key
	for z := 32; z < 126; z++ {
		results = append(results, &Result{
			key:   uint8(z),
			text:  deXor(w, uint8(z)),
			score: 0,
		})
	}

	// Score each attempt using relative frequencies
	for k, v := range results {
		dFreqs, total := mapMaker(v.text)
		score := float64(0)

		// Compare frequency distributions
		for char, dFreq := range dFreqs {
			if rFreq, exists := rFreqs[char]; exists {
				// Add up the differences in frequencies
				diff := dFreq - rFreq
				score += diff * diff // Square the difference to penalize larger deviations
			} else {
				// Penalize characters that don't appear in reference text
				score += dFreq * dFreq
			}
		}
		results[k].total = total
		results[k].score = score
	}

	sort.Sort(sortByScore{results})

	// Print top 10 results
	for i, v := range results[:10] {
		fmt.Printf("Text input size: %d\n", v.total)
		fmt.Printf("Result %d - Key: %c (0x%02x) Score: %f \n", i+1, v.key, v.key, v.score)
		fmt.Printf("Decrypted text: %s\n", string(v.text))
		fmt.Printf("====================================\n")
	}
}
