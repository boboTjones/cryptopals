/*
8. Detecting ECB

At the following URL are a bunch of hex-encoded ciphertexts:

  https://gist.github.com/3132928

One of them is ECB encrypted. Detect it.

Remember that the problem with ECB is that it is stateless and
deterministic; the same 16 byte plaintext block will always produce
the same 16 byte ciphertext.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"util"
)

var fileName string

type Line struct {
	Number int
	Score  int
}
type Lines []Line
type sortLinesByScore struct{ Lines }

func (line Lines) Len() int                   { return len(line) }
func (line Lines) Swap(i, j int)              { line[i], line[j] = line[j], line[i] }
func (s sortLinesByScore) Less(i, j int) bool { return s.Lines[i].Score > s.Lines[j].Score }

func twoStep(in []byte) []string {
	out := make([]string, 0)
	if len(in) > 0 {
		junk := bytes.Split(in, nil)
		x := len(junk)
		tmp := make([][]byte, 0)

		for i := 0; i < x; i += 2 {
			q := bytes.Join(junk[i:i+2], nil)
			tmp = append(tmp, q)
		}

		for _, v := range tmp {
			out = append(out, string(v))
		}
	}
	return out
}

func compare(c []string, n int) int {
	var ret int
	v := chunk(c, n)
	x := len(v) - 1
	for y := 0; y < x; y++ {
		hi := v[y]
		bi := v[y+1:]
		for _, b := range bi {
			for i, h := range hi {
				if h == b[i] {
					ret++
				}
			}
		}
	}
	return ret
}

func chunk(blob []string, csize int) [][]string {
	var fin = make([][]string, 0)
	x := 0
	for i := 0; i < len(blob); i += csize {
		fd := []string(blob[i:(x + csize)])
		fin = append(fin, fd)
		x += csize
	}
	return fin
}

func init() {
	flag.StringVar(&fileName, "f", fileName, "File to decrypt")
}

func main() {
	flag.Parse()

	data := bytes.NewBuffer(make([]byte, 0))

	if fileName != "" {
		d := util.SlurpFromFile(fileName)
		data.Write(d)
	} else {
		fmt.Println("need input")
		os.Exit(1)
	}

	src := bytes.Split(data.Bytes(), []byte("\n"))

	result := Lines{}

	for i, s := range src {
		if len(s) > 0 {
			str := twoStep(s)
			score := compare(str, 16)
			result = append(result, Line{Number: i, Score: score})
		}

	}

	sort.Sort(sortLinesByScore{result})
	fmt.Printf("The score for line %d is %d.\n", result[0].Number, result[0].Score)
}
