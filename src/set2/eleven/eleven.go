/*
Now that you have ECB and CBC working:

Write a function to generate a random AES key; that's just 16 random bytes.

Write a function that encrypts data under an unknown key --- that is, a function
that generates a random key and encrypts under it.

The function should look like:

encryption_oracle(your-input)
=> [MEANINGLESS JIBBER JABBER]
Under the hood, have the function append 5-10 bytes (count chosen randomly) before
the plaintext and 5-10 bytes after the plaintext.

Now, have the function choose to encrypt under ECB 1/2 the time, and under CBC the
other half (just use random IVs each time for CBC). Use rand(2) to decide which to use.

Detect the block cipher mode the function is using each time. You should end up with
a piece of code that, pointed at a block box that might be encrypting ECB or CBC,
tells you which one is happening.
*/

package main

import (
	"bytes"
	"cc/util"
	"fmt"
	"math/rand"
	"time"
)

const BlockSize = 16

func chance(src []byte) []byte {
	// random key
	key := util.RandString(16)
	iv := util.RandString(16)
	// random bytes.
	// Jesus. Really?
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	in := util.RandString(r1.Intn(8) + 2)
	s2 := rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s2)
	tail := util.RandString(r2.Intn(8) + 2)
	in = append(in, src...)
	in = append(in, tail...)
	// guess i gotta pad it, too
	if len(in)%BlockSize != 0 {
		pddr := new(util.Padder)
		pddr.BlockSize = BlockSize
		pddr.Data = &bytes.Buffer{}
		pddr.Data.Write(in)
		pddr.Padfoot()
		in = pddr.Data.Bytes()
	}
	out := make([]byte, len(in))
	cbc := util.NewCBC(iv, key)
	cbc.Encrypt(out, in)
	return out
}

func main() {
	//src := util.SlurpFromFile("/Users/erin/codebase/cryptochallenges/randomtxt.txt")
	src := make([]byte, 40)
	for i, _ := range src {
		src[i] = byte(0x41)
	}
	out := chance(src)
	fmt.Println(out)
}
