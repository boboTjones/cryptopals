package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

type CharSum struct {
	Char  byte
	Total int
}

type CharMap []CharSum
type sortByTotal struct{ CharMap }

func (c CharMap) Len() int               { return len(c) }
func (c CharMap) Swap(i, j int)          { c[i], c[j] = c[j], c[i] }
func (s sortByTotal) Less(i, j int) bool { return s.CharMap[i].Total > s.CharMap[j].Total }

// maybe put this in another util? maybe make this whole file a util?
func FetchStuff(t string) []byte {
	r, err := http.Get(t)
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

func Slurp(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
	}
	return data
}

func CharCount(text []byte) (m map[byte]int) {
	m = make(map[byte]int)
	for _, t := range text {
		m[t] += 1
	}
	return
}

func sortMap(in map[byte]int) CharMap {
	cmap := CharMap{}
	for c, t := range in {
		cmap = append(cmap, CharSum{c, t})
	}
	sort.Sort(sortByTotal{cmap})
	return cmap
}

func MapMaker(text []byte) CharMap {
	return sortMap(CharCount(text))
}

/*
func main() {
    cfm := mapMaker(fetchStuff("http://www.gutenberg.org/files/205/205-h/205-h.htm"))
    for _, v := range(cfm) {
        fmt.Println(v)
    }
}
*/
