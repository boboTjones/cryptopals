package main

import (
	"bytes"
	"cc/util"
	"flag"
	"fmt"
	"sort"
)

var fileName string

type Ham struct {
	Distance int
	Score    float32
	Keysize  int
}
type Hammed []Ham
type sortHamsByScore struct{ Hammed }
type sortHamsByDistance struct{ Hammed }

func (h Hammed) Len() int                       { return len(h) }
func (h Hammed) Swap(i, j int)                  { h[i], h[j] = h[j], h[i] }
func (s sortHamsByScore) Less(i, j int) bool    { return s.Hammed[i].Score < s.Hammed[j].Score }
func (s sortHamsByDistance) Less(i, j int) bool { return s.Hammed[i].Distance < s.Hammed[j].Distance }

type Block struct {
	Blob        []byte
	String      []byte
	Score       int
	Character   uint8
	BlockNumber int
}
type Blocks []Block
type sortBlocksByScore struct{ Blocks }
type sortBlocksByNumber struct{ Blocks }

func (b Blocks) Len() int                      { return len(b) }
func (b Blocks) Swap(i, j int)                 { b[i], b[j] = b[j], b[i] }
func (s sortBlocksByScore) Less(i, j int) bool { return s.Blocks[i].Score > s.Blocks[j].Score }
func (s sortBlocksByNumber) Less(i, j int) bool {
	return s.Blocks[i].BlockNumber < s.Blocks[j].BlockNumber
}

func HammingDistance(str1 []byte, str2 []byte) int {
	d := 0
	for i := 0; i < len(str1); i++ {
		for b := 0; b < 7; b++ {
			// if the bit is set, x will be zero
			x := str1[i] & (1 << uint(b))
			// if the bit is set, y will be zero
			y := str2[i] & (1 << uint(b))
			switch {
			case x == 0 && y > 0:
				d++
			case x > 0 && y == 0:
				d++
			}
		}
	}
	return d
}

func shuffle(blob []byte, csize int) [][]byte {
	var fin = make([][]byte, csize)
	for i := 0; i < len(blob); i++ {
		fin[i%csize] = append(fin[i%csize], blob[i])
	}
	return fin
}

func chunk(blob []byte, csize int) [][]byte {
	var fin = make([][]byte, 0)
	x := 0
	for i := 0; i < len(blob)/csize; i += csize {
		fin = append(fin, []byte(blob[i:(x+csize)]))
		x += csize
	}
	fin = append(fin, []byte(blob[x:]))
	return fin
}

func scoreText(in []byte) int {
	var res int
	es := []string{"U", "C", "D", "L", "H", "R", "S", "N", "I", "O", "A", "T", "E", "u", "l", "d", "r", "h", "s", "n", "i", "o", "a", "t", "e", " "}
	for _, v := range in {
		for i, x := range es {
			if string(v) == x {
				res += i
			}
		}
	}
	return res
}

func scoreAscii(in []byte) int {
	var res int
	for _, v := range in {
		if v > 31 && v < 127 {
			res++
		}
		if v == 0x0a || v == 0x0d {
			res++
		}
	}
	return res
}

func scoreSpace(in []byte) int {
	var res int
	for _, v := range in {
		if v == 32 {
			res++
		}

	}
	return res
}

func init() {
	flag.StringVar(&fileName, "f", fileName, "File to decrypt")
}

func main() {
	flag.Parse()

	data := bytes.NewBuffer(make([]byte, 0))

	blob := util.SlurpFromFile(fileName)
	plain := util.Decode64(string(blob))
	data.Write(plain)

	hammed := Hammed{}
	msg := data.Bytes()

	for i := 30; i > 1; i-- {
		sab := chunk(msg, i)
		d := 0
		x := 1
		for t := 0; t < 3; t++ {
			d += HammingDistance(sab[t], sab[t+1])
			x++
		}
		distance := d / (x)
		score := float32(distance) / float32(i)
		hammed = append(hammed, Ham{(distance), score, i})
	}

	sort.Sort(sortHamsByScore{hammed})

	// this is cheating. :-)
	keySize := hammed[3].Keysize
	sab := shuffle(msg, keySize)
	all := Blocks{}

	for i, b := range sab {
		shuffled := Blocks{}
		for k := 32; k < 127; k++ {
			block := Block{Blob: b, BlockNumber: i}
			block.Character = uint8(k)
			block.String = util.DeXor(b, block.Character)
			block.Score = scoreText(block.String)
			shuffled = append(shuffled, block)
		}
		sort.Sort(sortBlocksByScore{shuffled})
		all = append(all, shuffled[0])
	}

	key := make([]byte, 0)
	for _, a := range all {
		key = append(key, a.Character)
	}

	var result []byte

	d := data.Bytes()
	for i := 0; i < len(d); i++ {
		k := key[i%len(key)]
		result = append(result, (d[i] ^ k))
	}
	fmt.Println(string(result))
}
