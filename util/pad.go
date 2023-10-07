package util

import "bytes"

type Padder struct {
	Data      *bytes.Buffer
	BlockSize int
}

func NewPadder(bs int) *Padder {
	d := &bytes.Buffer{}
	return &Padder{Data: d, BlockSize: bs}
}

// It'd be nice if this was in-place, but I don't feel like re-writing it.

func (self *Padder) Padfoot() {
	var k int

	if self.Data.Len() > self.BlockSize {
		blocks := Chunk(self.Data.Bytes(), self.BlockSize)
		lb := blocks[len(blocks)-1:][0]
		k = len(lb)
	} else {
		k = self.Data.Len()

	}
	p := self.BlockSize % k
	if k < self.BlockSize {
		for i := k; i < self.BlockSize; i++ {
			self.Data.WriteByte(byte(p))
		}
	}

}
