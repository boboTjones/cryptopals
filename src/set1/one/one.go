/*
1. Convert hex to base64 and back.

The string:

  49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d

should produce:

  SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t
*/

package main

import (
  "encoding/base64"
  "encoding/hex"
  "fmt"
)

func main() {
  strstr := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
  foo := make([]byte, len(strstr)/2)
  _, err := hex.Decode(foo, []byte(strstr))
  if err != nil {
    fmt.Println("Something bad happened")
  }
  fmt.Println(base64.StdEncoding.EncodeToString(foo))
}
