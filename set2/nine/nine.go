/*
Implement PKCS#7 padding

A block cipher transforms a fixed-sized block (usually 8 or 16 bytes) of plaintext into ciphertext.
But we almost never want to transform a single block; we encrypt irregularly-sized messages.

One way we account for irregularly-sized messages is by padding, creating a plaintext that is an
even multiple of the blocksize. The most popular padding scheme is called PKCS#7.

So: pad any block to a specific block length, by appending the number of bytes of padding to the
end of the block. For instance,

"YELLOW SUBMARINE"

... padded to 20 bytes would be:

"YELLOW SUBMARINE\x04\x04\x04\x04"

*/

package main

import (
	"bytes"
	"fmt"
	"github.com/bobotjones/cryptopals/util"
)

type Padder struct {
	data      bytes.Buffer
	blockSize int
}

func (self *Padder) Pad() {
	k := self.data.Len()
	p := self.blockSize % k
	if k < self.blockSize {
		for i := k; i < self.blockSize; i++ {
			self.data.WriteByte(byte(p))
		}
	}

}

func (self *Padder) Pad2() {
	var k int

	if self.data.Len() > self.blockSize {
		blocks := util.Chunk(self.data.Bytes(), self.blockSize)
		lb := blocks[len(blocks)-1:][0]
		k = len(lb)
	} else {
		k = self.data.Len()

	}
	p := self.blockSize % k
	if k < self.blockSize {
		for i := k; i < self.blockSize; i++ {
			self.data.WriteByte(byte(p))
		}
	}
}

func main() {

	pddr := new(Padder)
	pddr.data.WriteString("YELLOW SUBMARINE")
	pddr.blockSize = 20
	pddr.Pad()
	fmt.Printf("%#v\n", pddr.data.String())
	pddr.Pad2()
	fmt.Printf("%#v\n", pddr.data.String())

}
