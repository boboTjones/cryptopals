package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"util"
)

type Profile struct {
	Email string `json:"email"`
	Uid   int    `json:"uid"`
	Role  string `json:"role"`
}

func qsParse(qs string) []uint8 {
	out := new(Profile)
	for _, v := range strings.Split(qs, "&") {
		tmp := strings.Split(v, "=")
		if len(tmp) != 2 {
			fmt.Printf("Found garbage %q\n", tmp)
			continue
		}
		switch tmp[0] {
		case "email":
			out.Email = tmp[1]
		case "uid":
			x, err := strconv.Atoi(tmp[1])
			if err != nil {
				panic(err)
			}
			out.Uid = x
		case "role":
			out.Role = tmp[1]
		default:
			fmt.Printf("Found garbage %q\n", tmp)
		}
	}
	ret, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	return ret

}

func profile_for(email string) string {
	var ret []string

	// eating & and =
	nopes := []string{"&", "="}
	for _, v := range nopes {
		email = strings.Replace(email, v, "", -1)
	}
	ret = append(ret, "email="+email)
	ret = append(ret, "uid=10&role=user")
	return strings.Join(ret, "&")
}

func encrypt(key, in []byte) []byte {
	ebc := util.NewECB(key)
	pddr := util.NewPadder(16)
	pddr.Data.Write(in)
	pddr.Padfoot()
	for i, v := range util.Chunk(pddr.Data.Bytes(), 16) {
		fmt.Printf("%d\t%q\n", i, v)
	}
	ret := make([]byte, pddr.Data.Len())
	ebc.Encrypt(ret, pddr.Data.Bytes())
	return ret
}

func decrypt(key, in []byte) []byte {
	ret := make([]byte, len(in))
	ecb := util.NewECB(key)
	ecb.Decrypt(ret, in)
	return ret
}

func main() {
	//qp := qsParse("foo=bar&baz=qux&zap=zazzle")
	//fmt.Printf("%q\n", qp)
	pq := profile_for("erin@f.bar.com")

	//  Generate a random AES key, then:
	key := util.RandString(16)
	//Encrypt the encoded user profile under the key; "provide" that to the "attacker".
	enc := encrypt(key, []byte(pq))

	dec := decrypt(key, enc)
	//Decrypt the encoded user profile and parse it.
	fmt.Printf("%q\n", qsParse(string(dec)))

	// Using only the user input to profile_for() (as an oracle to generate "valid"
	// ciphertexts) and the ciphertexts themselves, make a role=admin profile.
	fmt.Println("Attack!")
	a1 := profile_for("eri@f.bar.com")
	c1 := encrypt(key, []byte(a1))
	b1 := c1[0:32]
	a2 := profile_for("XXXXXXXXX") + "admin           "
	c2 := encrypt(key, []byte(a2))
	e := len(c2) - 16

	b2 := c2[e:]
	as := append(b1, b2...)
	o := decrypt(key, as)
	fmt.Printf("%q\n", qsParse(string(o)))
}
