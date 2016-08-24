package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func qsParse(qs string) []uint8 {
	out := make(map[string]interface{})
	for _, v := range strings.Split(qs, "&") {
		tmp := strings.Split(v, "=")
		out[tmp[0]] = tmp[1]
	}
	ret, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	return ret

}

func profile_for(email string) string {
	var ret []string
	// ew
	tmp := make(map[string]interface{})
	tmp["email"] = email
	tmp["uid"] = 10
	tmp["role"] = "user"
	out, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T\n", tmp)
	fmt.Printf("%T\n", out)

	// ew ew ew
	for k, v := range out {
		//	ret = append(ret, k+"="+v)
		fmt.Printf("%v\n", string(k))
		fmt.Printf("%v\n", string(v))
	}
	return strings.Join(ret, "&")
}

func main() {
	qp := qsParse("foo=bar&baz=qux&zap=zazzle")
	fmt.Printf("%T\n", qp)
	pq := profile_for("foo@bar.com")
	fmt.Printf("%s\n", pq)
}
