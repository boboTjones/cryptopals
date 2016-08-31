package main

import (
	"cc/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func qsParse(qs string) []uint8 {
	out := make(map[string]interface{})
	for _, v := range strings.Split(qs, "&") {
		tmp := strings.Split(v, "=")
		fmt.Println(tmp)
		if len(tmp) > 1 {
			out[tmp[0]] = tmp[1]
		}
	}
	ret, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	return ret

}

// Assuming this is leading up to using JSON?
func profile_for(email string) string {
	var ret []string
	// eating & and =
	// will cause a lookup to fail.
	nopes := []string{"&", "="}
	for _, v := range nopes {
		email = strings.Replace(email, v, "", -1)
	}

	// ew
	tmp := make(map[string]interface{})
	tmp["email"] = email
	tmp["uid"] = 10
	tmp["role"] = "user"

	// ew
	for k, v := range tmp {
		switch x := v.(type) {
		case int:
			ret = append(ret, k+"="+strconv.Itoa(x))
		case float64:
			ret = append(ret, k+"="+strconv.FormatFloat(x, 'f', 2, 32))
		default:
			ret = append(ret, k+"="+v.(string))
		}
	}
	return strings.Join(ret, "&")
}

//  Now, two more easy functions.

func encrypt(key, in []byte) []byte {
	// Because I'm lazy. XXX TODO(erin)
	cbc := util.NewCBC(key, key)
	pddr := util.NewPadder(16)
	pddr.Data.Write(in)
	pddr.Padfoot()
	ret := make([]byte, pddr.Data.Len())
	cbc.Encrypt(ret, pddr.Data.Bytes())
	return ret
}

func decrypt(key, in []byte) []byte {
	ret := make([]byte, len(in))
	// Because I'm lazy. XXX TODO(erin)
	cbc := util.NewCBC(key, key)
	cbc.Decrypt(ret, in)
	return ret
}

func main() {
	//qp := qsParse("foo=bar&baz=qux&zap=zazzle")
	//fmt.Printf("%q\n", qp)
	pq := profile_for("erin@f.bar.com")
	for _, v := range util.Chunk([]byte(pq), 16) {
		fmt.Printf("%q\n", v)
	}

	//  Generate a random AES key, then:
	key := util.RandString(16)
	//Encrypt the encoded user profile under the key; "provide" that to the "attacker".
	enc := encrypt(key, []byte(pq))

	dec := decrypt(key, enc)
	//Decrypt the encoded user profile and parse it.
	fmt.Printf("%q\n", qsParse(string(dec)))

	// Using only the user input to profile_for() (as an oracle to generate "valid"
	// ciphertexts) and the ciphertexts themselves, make a role=admin profile.
	// What?
	fmt.Println("Attack!")
	a1 := profile_for("erin@f.bar.com")
	c1 := encrypt(key, []byte(a1))
	b1 := c1[0:32]
	a2 := profile_for("er@bar.com") + "admin"
	c2 := encrypt(key, []byte(a2))
	b2 := c2[(len(c2) - 16):]
	as := append(b1, b2...)
	o := decrypt(key, as)
	fmt.Printf("%q\n", qsParse(string(o)))
}
